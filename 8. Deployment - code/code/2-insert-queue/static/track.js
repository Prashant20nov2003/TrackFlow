var _goTracker=(()=>{var c=class{id="";siteId="";referrer="";isTouch=!1;constructor(t,e){this.siteId=t,this.referrer=e,this.isTouch="ontouchstart"in window||navigator.maxTouchPoints>0;let r=this.getSession("id");r&&(this.id=r)}getSession(t){t=`__got_${t}__`;let e=localStorage.getItem(t);return e?JSON.parse(e):null}setSession(t,e){t=`__got_${t}__`,localStorage.setItem(t,JSON.stringify(e))}identify(t){this.id=t,this.setSession("id",t)}track(t,e){let r={tracking:{type:e=="Page views"?"page":"event",identity:this.id,ua:navigator.userAgent,event:t,category:e,referrer:this.referrer,isTouchDevice:this.isTouch},site_id:this.siteId};this.trackRequest(r)}page(t){this.track(t,"Page views")}trackRequest(t){let e=JSON.stringify(t),r=`http://localhost:9876/track?data=${btoa(e)}`,a=new Image;a.src=r}};((i,t)=>{let e=t.currentScript?.dataset;if(e){if(!e.siteid){console.error("you must have a data-siteid in your script tag.");return}}else{console.error("you must have a data-siteid in your script tag.");return}let r=i.location.pathname,a="",n=t.referrer;n&&n.indexOf(`${i.location.protocol}//${i.location.host}`)==0&&(a=n);let s=new c(e.siteid,a);i._got=i._got||s,s.page(r);let o=window.history;if(o.pushState){let g=o.pushState;o.pushState=function(){g.apply(this,arguments),s.page(i.location.pathname)},window.addEventListener("popstate",()=>{s.page(i.location.pathname)})}i.addEventListener("hashchange",()=>{s.page(t.location.hash)},!1)})(window,document);})();
