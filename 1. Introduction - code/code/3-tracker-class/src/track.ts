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
  private trackRequest(payload: TrackPayload) {
    const s = JSON.stringify(payload);
    const url = `http://localhost:9876/track?data=${btoa(s)}`;

    const img = new Image();
    img.src = url;
  }
}
((w, d) => {
  w._got = new Tracker();
  console.log("tracker loaded");
})(window, document);
