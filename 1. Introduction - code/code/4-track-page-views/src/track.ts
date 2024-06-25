interface TrackingData {
  type: "event" | "page";
  identity: string;
  ua: string;
  event: string;
  category: string;
  referrer: string;
}

interface TrackPayload {
  tracking: TrackingData;
  site_id: string;
}

class Tracker {
  track(event: string, category: string) {
    const payload: TrackPayload = {
      tracking: {
        type: "event",
        identity: "todo-id",
        ua: navigator.userAgent,
        event: event,
        category: category,
        referrer: "todo-referrer",
      },
      site_id: "todo-site-id",
    };
    this.trackRequest(payload);
  }

  page(path: string) {
    this.track(path, "Page views");
  }
  private trackRequest(payload: TrackPayload) {
    const s = JSON.stringify(payload);
    const url = `http://localhost:9876/track?data=${btoa(s)}`;

    const img = new Image();
    img.src = url;
  }
}
((w, d) => {
  const path = w.location.pathname;

  let tracker = new Tracker();

  w._got = w._got || tracker;

  tracker.page(path);

  const his = window.history;
  if (his.pushState) {
    const originalFn = his["pushState"];
    his.pushState = function () {
      originalFn.apply(this, arguments);
      tracker.page(w.location.pathname);
    };

    window.addEventListener("popstate", () => {
      tracker.page(w.location.pathname);
    });
  }

  w.addEventListener(
    "hashchange",
    () => {
      tracker.page(d.location.hash);
    },
    false
  );
})(window, document);
