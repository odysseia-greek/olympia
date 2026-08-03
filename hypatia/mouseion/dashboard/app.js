const number = new Intl.NumberFormat();

async function getJSON(path) {
  const response = await fetch(path, { cache: "no-store" });
  if (!response.ok) throw new Error(`${response.status} ${response.statusText}`);
  return response.json();
}

function text(id, value) {
  document.getElementById(id).textContent = value;
}

function short(value, length = 14) {
  if (!value) return "—";
  return value.length > length ? `${value.slice(0, length)}…` : value;
}

function statusClass(status) {
  if (status >= 500) return "status error";
  if (status >= 400) return "status warn";
  return "status";
}

function renderEvents(events) {
  const rows = document.getElementById("eventRows");
  rows.replaceChildren();
  if (!events.length) {
    const row = rows.insertRow();
    const cell = row.insertCell();
    cell.colSpan = 5;
    cell.className = "empty";
    cell.textContent = "No requests observed yet.";
    return;
  }
  for (const event of events) {
    const row = rows.insertRow();
    const timestamp = event.timestamp ? new Date(event.timestamp) : null;
    row.insertCell().textContent = timestamp && !Number.isNaN(timestamp.valueOf()) ? timestamp.toLocaleTimeString() : "—";

    const request = row.insertCell();
    const method = document.createElement("span");
    method.className = "method";
    method.textContent = event.method || "—";
    request.append(method, document.createTextNode(event.path || "—"));

    const status = document.createElement("span");
    status.className = statusClass(event.status);
    status.textContent = event.status || "—";
    row.insertCell().append(status);

    const session = document.createElement("span");
    session.className = "session";
    session.title = event.sessionId || "";
    session.textContent = short(event.sessionId);
    row.insertCell().append(session);

    const trace = document.createElement("span");
    trace.className = "trace";
    trace.title = event.traceId || "";
    trace.textContent = short(event.traceId);
    row.insertCell().append(trace);
  }
}

function renderRanking(id, values, emptyLabel) {
  const list = document.getElementById(id);
  list.replaceChildren();
  if (!values.length) {
    const item = document.createElement("li");
    item.className = "empty";
    item.textContent = emptyLabel;
    list.append(item);
    return;
  }
  for (const entry of values) {
    const item = document.createElement("li");
    const value = document.createElement("span");
    value.className = "value";
    value.title = entry.value;
    value.textContent = entry.value;
    const count = document.createElement("span");
    count.className = "count";
    count.textContent = number.format(entry.count);
    item.append(value, count);
    list.append(item);
  }
}

async function refresh() {
  const button = document.getElementById("refresh");
  button.disabled = true;
  try {
    const [summary, events, paths, sessions] = await Promise.all([
      getJSON("/api/summary"),
      getJSON("/api/events?limit=50"),
      getJSON("/api/paths?limit=8"),
      getJSON("/api/sessions?limit=8"),
    ]);
    text("events", number.format(summary.events));
    text("sessions", number.format(summary.sessions));
    text("paths", number.format(summary.paths));
    text("traced", number.format(summary.tracedEvents));
    text("updated", new Date(summary.generatedAt).toLocaleTimeString());
    renderEvents(events);
    renderRanking("pathRows", paths, "No path data yet.");
    renderRanking("sessionRows", sessions, "No session data yet.");
  } catch (error) {
    const rows = document.getElementById("eventRows");
    rows.innerHTML = "";
    const row = rows.insertRow();
    const cell = row.insertCell();
    cell.colSpan = 5;
    cell.className = "empty error-message";
    cell.textContent = `Dashboard unavailable: ${error.message}`;
  } finally {
    button.disabled = false;
  }
}

document.getElementById("refresh").addEventListener("click", refresh);
refresh();
setInterval(refresh, 15_000);
