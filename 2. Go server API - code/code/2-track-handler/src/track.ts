interface TrackingData {
  type: "event" | "page";
  identity: string;
  ua: string;
  event: string;
  category: string;
  referrer: string;
  isTouchDevice: boolean;
}

interface TrackPayload {
  tracking: TrackingData;
  site_id: string;
}

class Tracker {
  private id: string = "";
  private siteId: string = "";
  private referrer: string = "";
  private isTouch = false;

  constructor(siteId: string, ref: string) {
    this.siteId = siteId;
    this.referrer = ref;
    this.isTouch = "ontouchstart" in window || navigator.maxTouchPoints > 0;

    const customId = this.getSession("id");
    if (customId) {
      this.id = customId;
    }
  }

  private getSession(key) {
    key = `__got_${key}__`;

    const s = localStorage.getItem(key);
    if (!s) return null;
    return JSON.parse(s);
  }

  private setSession(key: string, value: any) {
    key = `__got_${key}__`;

    localStorage.setItem(key, JSON.stringify(value));
  }

  identify(customId: string) {
    this.id = customId;
    this.setSession("id", customId);
  }

  track(event: string, category: string) {
    const payload: TrackPayload = {
      tracking: {
        type: category == "Page views" ? "page" : "event",
        identity: this.id,
        ua: navigator.userAgent,
        event: event,
        category: category,
        referrer: this.referrer,
        isTouchDevice: this.isTouch,
      },
      site_id: this.siteId,
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
  const ds = d.currentScript?.dataset;
  if (!ds) {
    console.error("you must have a data-siteid in your script tag.");
    return;
  } else if (!ds.siteid) {
    console.error("you must have a data-siteid in your script tag.");
    return;
  }

  const path = w.location.pathname;

  let externalReferrer = "";
  const ref = d.referrer;
  if (ref && ref.indexOf(`${w.location.protocol}//${w.location.host}`) == 0) {
    externalReferrer = ref;
  }

  let tracker = new Tracker(ds.siteid, externalReferrer);

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
