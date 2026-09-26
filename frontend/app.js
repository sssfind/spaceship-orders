(() => {
  const API_BASE = (() => {
    const fromQuery = new URLSearchParams(window.location.search).get("api");
    if (fromQuery) return fromQuery.replace(/\/$/, "");
    return "";
  })();

  const SESSION_KEY = "spaceship_session_uuid";

  const els = {
    authPanel: document.getElementById("auth-panel"),
    authForm: document.getElementById("auth-form"),
    authLogin: document.getElementById("auth-login"),
    authPassword: document.getElementById("auth-password"),
    authStatus: document.getElementById("auth-status"),
    authBar: document.getElementById("auth-bar"),
    authUser: document.getElementById("auth-user"),
    logoutBtn: document.getElementById("logout-btn"),
    appPanel: document.getElementById("app-panel"),
    form: document.getElementById("create-form"),
    createPartForm: document.getElementById("create-part-form"),
    partsBody: document.getElementById("parts-body"),
    partsStatus: document.getElementById("parts-status"),
    partsRefresh: document.getElementById("parts-refresh-btn"),
    tbody: document.getElementById("orders-body"),
    status: document.getElementById("status"),
    refresh: document.getElementById("refresh-btn"),
    payDialog: document.getElementById("pay-dialog"),
    payForm: document.getElementById("pay-form"),
    payOrderId: document.getElementById("pay-order-id"),
    paymentMethod: document.getElementById("payment-method"),
  };

  let payTargetUuid = null;

  function getSession() {
    return localStorage.getItem(SESSION_KEY) || "";
  }

  function setSession(uuid) {
    if (uuid) localStorage.setItem(SESSION_KEY, uuid);
    else localStorage.removeItem(SESSION_KEY);
  }

  function setAuthStatus(message, kind = "") {
    els.authStatus.textContent = message || "";
    els.authStatus.className = "status" + (kind ? ` status--${kind}` : "");
  }

  function setStatus(message, kind = "") {
    els.status.textContent = message || "";
    els.status.className = "status" + (kind ? ` status--${kind}` : "");
  }

  function setPartsStatus(message, kind = "") {
    els.partsStatus.textContent = message || "";
    els.partsStatus.className = "status" + (kind ? ` status--${kind}` : "");
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
    const session = getSession();
    const res = await fetch(`${API_BASE}${path}`, {
      headers: {
        Accept: "application/json",
        ...(options.body ? { "Content-Type": "application/json" } : {}),
        ...(session ? { "X-Session-Uuid": session } : {}),
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

  function showLoggedOut() {
    els.authPanel.classList.remove("hidden");
    els.appPanel.classList.add("hidden");
    els.authBar.classList.add("hidden");
    els.authUser.textContent = "";
  }

  function showLoggedIn(user) {
    els.authPanel.classList.add("hidden");
    els.appPanel.classList.remove("hidden");
    els.authBar.classList.remove("hidden");
    els.authUser.textContent = user?.login
      ? `${user.login} (${shortUuid(user.user_uuid)})`
      : "Сессия активна";
  }

  async function refreshSessionUI() {
    const session = getSession();
    if (!session) {
      showLoggedOut();
      return false;
    }
    try {
      const { data } = await api("/api/v1/auth/whoami");
      showLoggedIn(data);
      return true;
    } catch {
      setSession("");
      showLoggedOut();
      return false;
    }
  }

  function selectedPartUuids() {
    return [...els.partsBody.querySelectorAll('input[type="checkbox"][data-part]:checked')].map(
      (el) => el.dataset.part
    );
  }

  function renderParts(parts) {
    if (!parts || parts.length === 0) {
      els.partsBody.innerHTML = `<tr><td colspan="6" class="empty">Каталог пуст</td></tr>`;
      return;
    }

    els.partsBody.innerHTML = parts
      .map((p) => {
        const disabled = Number(p.in_stock) <= 0 ? "disabled" : "";
        return `
          <tr>
            <td><input type="checkbox" data-part="${p.part_uuid}" ${disabled} /></td>
            <td>${p.name}</td>
            <td>${p.category}</td>
            <td>${Number(p.price).toFixed(2)}</td>
            <td>${p.in_stock}</td>
            <td class="mono" title="${p.part_uuid}">${shortUuid(p.part_uuid)}</td>
          </tr>`;
      })
      .join("");
  }

  async function loadParts() {
    setPartsStatus("Загрузка каталога…");
    try {
      const { data } = await api("/api/v1/parts");
      renderParts(Array.isArray(data) ? data : []);
      setPartsStatus(`Деталей: ${Array.isArray(data) ? data.length : 0}`, "ok");
    } catch (err) {
      if (err.status === 401) {
        setSession("");
        showLoggedOut();
        setAuthStatus("Сессия истекла, войдите снова", "err");
        return;
      }
      els.partsBody.innerHTML = `<tr><td colspan="6" class="empty">Не удалось загрузить каталог</td></tr>`;
      setPartsStatus(err.message, "err");
    }
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
      if (err.status === 401) {
        setSession("");
        showLoggedOut();
        setAuthStatus("Сессия истекла, войдите снова", "err");
        return;
      }
      els.tbody.innerHTML = `<tr><td colspan="5" class="empty">Не удалось загрузить заказы</td></tr>`;
      setStatus(err.message, "err");
    }
  }

  async function loadAppData() {
    await Promise.all([loadParts(), loadOrders()]);
  }

  els.authForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const action = e.submitter?.value || "login";
    const login = els.authLogin.value.trim();
    const password = els.authPassword.value;
    if (!login || !password) {
      setAuthStatus("Укажите логин и пароль", "err");
      return;
    }

    const path = action === "register" ? "/api/v1/auth/register" : "/api/v1/auth/login";
    setAuthStatus(action === "register" ? "Регистрация…" : "Вход…");
    try {
      if (action === "register") {
        await api(path, {
          method: "POST",
          body: JSON.stringify({ login, password }),
        });
        setAuthStatus("Аккаунт создан, выполняем вход…", "ok");
      }

      const { data } = await api("/api/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({ login, password }),
      });
      setSession(data.session_uuid);
      els.authPassword.value = "";
      setAuthStatus("");
      showLoggedIn(data.user);
      await loadAppData();
    } catch (err) {
      setAuthStatus(err.message, "err");
    }
  });

  els.logoutBtn.addEventListener("click", async () => {
    try {
      await api("/api/v1/auth/logout", { method: "POST" });
    } catch {
      // ignore logout errors
    }
    setSession("");
    showLoggedOut();
    setAuthStatus("Вы вышли", "ok");
  });

  els.form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const part_uuids = selectedPartUuids();
    if (part_uuids.length === 0) {
      setStatus("Выберите хотя бы одну деталь в каталоге", "err");
      return;
    }

    setStatus("Создание заказа…");
    try {
      const { data } = await api("/api/v1/orders", {
        method: "POST",
        body: JSON.stringify({ part_uuids }),
      });
      setStatus(`Создан заказ ${data.order_uuid} (цена ${data.total_price})`, "ok");
      await loadAppData();
    } catch (err) {
      setStatus(err.message, "err");
    }
  });

  els.createPartForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const name = document.getElementById("part-name").value.trim();
    const price = Number(document.getElementById("part-price").value);
    const category = document.getElementById("part-category").value.trim() || "GENERAL";
    const in_stock = Number(document.getElementById("part-stock").value);

    setPartsStatus("Создание детали…");
    try {
      await api("/api/v1/parts", {
        method: "POST",
        body: JSON.stringify({ name, price, category, in_stock }),
      });
      els.createPartForm.reset();
      document.getElementById("part-category").value = "GENERAL";
      document.getElementById("part-stock").value = "10";
      setPartsStatus("Деталь добавлена", "ok");
      await loadParts();
    } catch (err) {
      setPartsStatus(err.message, "err");
    }
  });

  els.refresh.addEventListener("click", () => loadOrders());
  els.partsRefresh.addEventListener("click", () => loadParts());

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
      setStatus(`Оплачено, payment ${data.transaction_uuid}`, "ok");
      await loadOrders();
    } catch (err) {
      setStatus(err.message, "err");
    }
  });

  (async () => {
    const ok = await refreshSessionUI();
    if (ok) await loadAppData();
  })();
})();
