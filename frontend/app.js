(() => {
  const API_BASE = (() => {
    const fromQuery = new URLSearchParams(window.location.search).get("api");
    if (fromQuery) return fromQuery.replace(/\/$/, "");
    // same origin when served by Order service
    return "";
  })();

  const els = {
    form: document.getElementById("create-form"),
    userUuid: document.getElementById("user-uuid"),
    partUuids: document.getElementById("part-uuids"),
    tbody: document.getElementById("orders-body"),
    status: document.getElementById("status"),
    refresh: document.getElementById("refresh-btn"),
    payDialog: document.getElementById("pay-dialog"),
    payForm: document.getElementById("pay-form"),
    payOrderId: document.getElementById("pay-order-id"),
    paymentMethod: document.getElementById("payment-method"),
  };

  let payTargetUuid = null;

  function setStatus(message, kind = "") {
    els.status.textContent = message || "";
    els.status.className = "status" + (kind ? ` status--${kind}` : "");
  }

  function parsePartUuids(raw) {
    return raw
      .split(/[\s,;]+/)
      .map((s) => s.trim())
      .filter(Boolean);
  }

  function shortUuid(uuid) {
    if (!uuid) return "—";
    return uuid.length > 13 ? `${uuid.slice(0, 8)}…${uuid.slice(-4)}` : uuid;
  }

  function statusBadge(status) {
    const key = String(status || "").toLowerCase().replace(/_/g, "-");
    const cls = ["pending-payment", "paid", "cancelled", "assembled", "completed"].includes(key)
      ? `badge--${key === "pending-payment" ? "pending" : key}`
      : "";
    return `<span class="badge ${cls}">${status || "—"}</span>`;
  }

  async function api(path, options = {}) {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: {
        Accept: "application/json",
        ...(options.body ? { "Content-Type": "application/json" } : {}),
        ...options.headers,
      },
      ...options,
    });

    const text = await res.text();
    let data = null;
    if (text) {
      try {
        data = JSON.parse(text);
      } catch {
        data = text;
      }
    }

    if (!res.ok) {
      const msg =
        (data && typeof data === "object" && data.message) ||
        (typeof data === "string" && data) ||
        `HTTP ${res.status}`;
      const err = new Error(msg);
      err.status = res.status;
      err.data = data;
      throw err;
    }

    return { status: res.status, data };
  }

  function renderOrders(orders) {
    if (!orders || orders.length === 0) {
      els.tbody.innerHTML = `<tr><td colspan="5" class="empty">Заказов пока нет</td></tr>`;
      return;
    }

    els.tbody.innerHTML = orders
      .map((o) => {
        const pending = o.status === "PENDING_PAYMENT";
        return `
          <tr data-uuid="${o.order_uuid}">
            <td class="mono" title="${o.order_uuid}">${shortUuid(o.order_uuid)}</td>
            <td class="mono" title="${o.user_uuid}">${shortUuid(o.user_uuid)}</td>
            <td>${Number(o.total_price).toFixed(2)}</td>
            <td>${statusBadge(o.status)}</td>
            <td>
              <div class="actions">
                <button type="button" class="btn btn--sm btn--primary" data-action="pay" ${pending ? "" : "disabled"}>Pay</button>
                <button type="button" class="btn btn--sm btn--ghost" data-action="cancel" ${pending ? "" : "disabled"}>Cancel</button>
                <button type="button" class="btn btn--sm btn--danger" data-action="delete">Delete</button>
              </div>
            </td>
          </tr>`;
      })
      .join("");
  }

  async function loadOrders() {
    setStatus("Загрузка списка…");
    try {
      const { data } = await api("/api/v1/orders");
      renderOrders(Array.isArray(data) ? data : []);
      setStatus(`Заказов: ${Array.isArray(data) ? data.length : 0}`, "ok");
    } catch (err) {
      els.tbody.innerHTML = `<tr><td colspan="5" class="empty">Не удалось загрузить заказы</td></tr>`;
      setStatus(err.message, "err");
    }
  }

  els.form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const user_uuid = els.userUuid.value.trim();
    const part_uuids = parsePartUuids(els.partUuids.value);
    if (!user_uuid || part_uuids.length === 0) {
      setStatus("Укажите user_uuid и хотя бы один part uuid", "err");
      return;
    }

    setStatus("Создание заказа…");
    try {
      const { data } = await api("/api/v1/orders", {
        method: "POST",
        body: JSON.stringify({ user_uuid, part_uuids }),
      });
      setStatus(`Создан заказ ${data.order_uuid} (цена ${data.total_price})`, "ok");
      await loadOrders();
    } catch (err) {
      setStatus(err.message, "err");
    }
  });

  els.refresh.addEventListener("click", () => loadOrders());

  els.tbody.addEventListener("click", async (e) => {
    const btn = e.target.closest("button[data-action]");
    if (!btn || btn.disabled) return;

    const row = btn.closest("tr[data-uuid]");
    const uuid = row?.dataset.uuid;
    if (!uuid) return;

    const action = btn.dataset.action;

    if (action === "pay") {
      payTargetUuid = uuid;
      els.payOrderId.textContent = uuid;
      els.payDialog.showModal();
      return;
    }

    if (action === "cancel") {
      if (!confirm(`Отменить заказ ${uuid}?`)) return;
      setStatus("Отмена…");
      try {
        await api(`/api/v1/orders/${uuid}/cancel`, { method: "POST" });
        setStatus("Заказ отменён", "ok");
        await loadOrders();
      } catch (err) {
        setStatus(err.message, "err");
      }
      return;
    }

    if (action === "delete") {
      if (!confirm(`Удалить заказ ${uuid}?`)) return;
      setStatus("Удаление…");
      try {
        await api(`/api/v1/orders/${uuid}`, { method: "DELETE" });
        setStatus("Заказ удалён", "ok");
        await loadOrders();
      } catch (err) {
        setStatus(err.message, "err");
      }
    }
  });

  els.payForm.addEventListener("submit", async (e) => {
    const submitter = e.submitter;
    if (!submitter || submitter.value !== "confirm") {
      payTargetUuid = null;
      return;
    }
    e.preventDefault();
    const uuid = payTargetUuid;
    const payment_method = els.paymentMethod.value;
    els.payDialog.close();
    payTargetUuid = null;
    if (!uuid) return;

    setStatus("Оплата…");
    try {
      const { data } = await api(`/api/v1/orders/${uuid}/pay`, {
        method: "POST",
        body: JSON.stringify({ payment_method }),
      });
      setStatus(`Оплачено, transaction ${data.transaction_uuid}`, "ok");
      await loadOrders();
    } catch (err) {
      setStatus(err.message, "err");
    }
  });

  loadOrders();
})();
