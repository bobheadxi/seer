(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["team"],{"0767":function(e,t,a){"use strict";a.r(t);var n=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"home"},[n("img",{attrs:{alt:"Vue logo",src:a("cf05")}}),n("h1",[e._v("Team "+e._s(e.teamID))]),n("p",[e._v("These stats are collected only from games where at least 4 members from this team played together.")]),e.team&&!e.loading?n("div",[n("div",[n("button",{on:{click:function(t){return e.copyMembersToClipboard()}}},[e._v("copy to clipboard")]),n("a",{attrs:{href:"http://na.op.gg/multi/query="+e.memberNames(),target:"_blank"}},[n("button",[e._v("open in na.op.gg")])]),n("button",{on:{click:function(t){return e.forceFetchTeam()}}},[e._v("refresh")])]),n("Overview",{attrs:{teamID:e.teamID}}),n("br"),n("Matches",{attrs:{teamID:e.teamID}}),n("br"),e.updateTriggered?n("div",[e._v("\n      Matches sync queued\n    ")]):e._e(),e.updateTriggered?e._e():n("button",{on:{click:function(t){return e.syncMatches()}}},[e._v("\n      Sync Matches\n    ")])],1):e._e(),e.error.occured?n("div",[e._v("\n    Oops an error occured! "+e._s(e.error.details)+"\n  ")]):e._e(),e.loading?n("div",[e._v("\n    Loading...\n  ")]):e._e()])},r=[],i=(a("7f7f"),a("96cf"),a("3b8d")),s=a("d225"),o=a("b0b4"),c=a("308d"),u=a("6bb5"),l=a("4e2b"),p=a("9ab4"),m=a("60a3"),v=a("4bb5"),d=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"matches"},[a("h2",[e._v("Matches")]),a("p",[e._v("This section is busted right now, don't bother looking.")]),e.matches?e._e():a("div",[e._v("\n    No matches found for this team.\n  ")]),e.matches?a("div",[a("p",[e._v(e._s(e.matches.length)+" matches retrieved.")]),e._l(e.matches,function(t){return a("div",{key:t.gameId},[a("h3",[e._v("Game "+e._s(t.details.gameId))]),a("p",[e._v(e._s(e.dateString(t.details.gameCreation)))]),e._l(t.details.participants,function(t){return a("div",{key:t.participantId},[e._v("\n        Participant: "+e._s(t.participantId)+"\n        "),a("br"),a("div",[a("img",{attrs:{src:e.champIcon(t.championId)}})]),a("br"),a("div",[e._v("\n          Spells:\n          "+e._s(e.spell(t.spell1Id).name)+"\n          "+e._s(e.spell(t.spell2Id).name)+"\n        ")]),a("br"),a("div",[e._v("\n          First Item ("+e._s(t.stats.item0)+"):\n          "),e.item(t.stats.item0)?e._e():a("div",[e._v("\n            No Item\n          ")]),e.item(t.stats.item0)?a("div",[e._v("\n            "+e._s(e.item(t.stats.item0).name)+"\n            "),a("img",{attrs:{src:e.itemIcon(t.stats.item0)}})]):e._e()]),a("br"),a("div",[a("div",[e._v("\n            Primary Perk ("+e._s(t.stats.perkPrimaryStyle)+"):\n            "+e._s(e.runes(t.stats.perkPrimaryStyle).key)+"\n          ")]),a("div",[e._v("\n            Secondary Perk ("+e._s(t.stats.perkSubStyle)+"):\n            "+e._s(e.runes(t.stats.perkSubStyle).key)+"\n          ")])]),a("hr")])})],2)})],2):e._e()])},h=[],b=a("0613"),f=a("8676"),g=a("349a"),y={namespace:b["a"].LEAGUE},_={namespace:b["a"].TEAMS},O=function(e){function t(){return Object(s["a"])(this,t),Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments))}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"dateString",value:function(e){return new Date(e).toLocaleString("en-US")}},{key:"matches",get:function(){return this.matchesData(this.teamID)}}]),t}(m["c"]);p["a"]([Object(m["b"])()],O.prototype,"teamID",void 0),p["a"]([Object(v["b"])(f["b"].MATCHES,_)],O.prototype,"matchesData",void 0),p["a"]([Object(v["b"])(g["b"].ITEM,y)],O.prototype,"item",void 0),p["a"]([Object(v["b"])(g["b"].ITEM_ICON,y)],O.prototype,"itemIcon",void 0),p["a"]([Object(v["b"])(g["b"].CHAMP,y)],O.prototype,"champ",void 0),p["a"]([Object(v["b"])(g["b"].CHAMP_ICON,y)],O.prototype,"champIcon",void 0),p["a"]([Object(v["b"])(g["b"].RUNES,y)],O.prototype,"runes",void 0),p["a"]([Object(v["b"])(g["b"].SPELL,y)],O.prototype,"spell",void 0),O=p["a"]([m["a"]],O);var I=O,j=I,k=a("2877"),T=Object(k["a"])(j,d,h,!1,null,null,null),D=T.exports,E=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"overview"},[a("h2",[e._v("Overview")]),a("div",e._l(e.team.members,function(t){return a("div",{key:t.id},[a("h3",[e._v("\n        "+e._s(t.name)+"\n        ("+e._s(e.playerOverviews&&e.playerOverviews[t.name]&&e.playerOverviews[t.name].tier?e.playerOverviews[t.name].tier+", ":"")+"\n        lv"+e._s(t.summonerLevel)+")\n        "),a("a",{attrs:{href:"https://na.op.gg/summoner/userName="+t.name,target:"_blank"}},[a("img",{attrs:{width:"16",height:"16",src:"https://lh3.googleusercontent.com/UdvXlkugn0bJcwiDkqHKG5IElodmv-oL4kHlNAklSA2sdlVWhojsZKaPE-qFPueiZg"}})])]),e.playerOverviews&&e.playerOverviews[t.name]&&e.playerOverviews[t.name].aggs?a("div",[a("div",[a("h5",[e._v("Most played lane and role")]),e._v("\n          "+e._s(e.playerOverviews[t.name].aggs.favourite.lane)+"\n          ("+e._s(e.playerOverviews[t.name].aggs.favourite.role)+")\n        ")]),a("div",[a("h5",[e._v("Most played champions")]),e._l(e.playerOverviews[t.name].aggs.favourite.champs,function(n){return a("img",{key:"fav-"+t.name+"-"+n,attrs:{src:e.champIcon(n)}})})],2),a("div",[a("h5",[e._v("Average stats")]),e._v("\n          "+e._s(e.playerOverviews[t.name].aggs.avg)+"\n        ")])]):e._e()])}),0)])},w=[],M=(a("7514"),a("ac6a"),a("456d"),a("55dd"),{namespace:b["a"].LEAGUE}),A={namespace:b["a"].TEAMS};function S(){return{vision:[],cs:[],dealt:[],taken:[],gold:[],champs:[],lanes:[],roles:[]}}function C(e,t){e.vision.push(t.stats.visionScore),e.cs.push(t.stats.totalMinionsKilled),e.dealt.push(t.stats.totalDamageDealtToChampions),e.taken.push(t.stats.totalDamageTaken),e.gold.push(t.stats.goldEarned),e.champs.push(t.championId),e.lanes.push(t.timeline.lane),e.roles.push(t.timeline.role)}function N(e){var t=e.reduce(function(e,t){return e+t},0);return(t/e.length).toFixed(2)}function P(e,t){var a=e.reduce(function(e,t){return e[t]=(e[t]||0)+1,e},{}),n=Object.keys(a).sort(function(e,t){return a[t]-a[e]});return t?n.slice(0,t):n[0]}function x(e){return{avg:{vision:N(e.vision),cs:N(e.cs),dealt:N(e.dealt),taken:N(e.taken),gold:N(e.gold)},favourite:{champs:P(e.champs,5),lane:P(e.lanes),role:P(e.roles)}}}var L=function(e){function t(){return Object(s["a"])(this,t),Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments))}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"idToName",value:function(){var e={};return this.teamData(this.teamID).members.forEach(function(t){e[t.accountId]=t.name}),e}},{key:"team",get:function(){return this.teamData(this.teamID)}},{key:"playerOverviews",get:function(){var e=this.matchesData(this.teamID);if(!e)return{};var t=this.idToName();console.debug("generating overviews",{matches:e.length});var a={},n={};return e.forEach(function(e){e.details.participantIdentities.forEach(function(r){var i=r.player.accountId;r.player.currentAccountId&&(i=r.player.currentAccountId);var s=t[i];if(s){var o=e.details.participants.find(function(e){return e.participantId===r.participantId});o?(n[s]||(n[s]={}),n[s].tier=o.highestAchievedSeasonTier,a[s]||(a[s]=S()),C(a[s],o)):console.debug("could not find participant ".concat(r.participantId))}})}),Object.keys(a).forEach(function(e){n[e].aggs=x(a[e])}),console.log("generated overviews",n),n}}]),t}(m["c"]);p["a"]([Object(m["b"])()],L.prototype,"teamID",void 0),p["a"]([Object(v["b"])(f["b"].MATCHES,A)],L.prototype,"matchesData",void 0),p["a"]([Object(v["b"])(f["b"].TEAM,{namespace:b["a"].TEAMS})],L.prototype,"teamData",void 0),p["a"]([Object(v["b"])(g["b"].ITEM,M)],L.prototype,"item",void 0),p["a"]([Object(v["b"])(g["b"].ITEM_ICON,M)],L.prototype,"itemIcon",void 0),p["a"]([Object(v["b"])(g["b"].CHAMP,M)],L.prototype,"champ",void 0),p["a"]([Object(v["b"])(g["b"].CHAMP_ICON,M)],L.prototype,"champIcon",void 0),p["a"]([Object(v["b"])(g["b"].RUNES,M)],L.prototype,"runes",void 0),p["a"]([Object(v["b"])(g["b"].SPELL,M)],L.prototype,"spell",void 0),L=p["a"]([m["a"]],L);var H=L,U=H,F=Object(k["a"])(U,E,w,!1,null,null,null),R=F.exports;function G(e){var t=document.createElement("textarea");t.value=e,t.setAttribute("readonly",""),document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t)}var q={namespace:b["a"].TEAMS},$={namespace:b["a"].LEAGUE},J=function(e){function t(){var e;return Object(s["a"])(this,t),e=Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments)),e.error={occured:!1},e.loading=!0,e.updateTriggered=!1,e}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"mounted",value:function(){var e=Object(i["a"])(regeneratorRuntime.mark(function e(){return regeneratorRuntime.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,this.fetchTeam({teamID:this.teamID});case 3:return e.next=5,this.fetchLeagueData({});case 5:e.next=10;break;case 7:e.prev=7,e.t0=e["catch"](0),this.error={occured:!0,details:e.t0};case 10:this.loading=!1;case 11:case"end":return e.stop()}},e,this,[[0,7]])}));function t(){return e.apply(this,arguments)}return t}()},{key:"syncMatches",value:function(){this.updateTriggered=!0,this.error={occured:!1};try{this.updateTeam({teamID:this.teamID})}catch(e){this.error={occured:!0,details:e},this.updateTriggered=!1}}},{key:"memberNames",value:function(){return this.team?this.team.members.map(function(e){return e.name}).join(","):""}},{key:"copyMembersToClipboard",value:function(){var e=this.memberNames();G(e)}},{key:"forceFetchTeam",value:function(){var e=Object(i["a"])(regeneratorRuntime.mark(function e(){return regeneratorRuntime.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return this.loading=!0,e.prev=1,e.next=4,this.fetchTeam({teamID:this.teamID,force:!0});case 4:e.next=9;break;case 6:e.prev=6,e.t0=e["catch"](1),this.error={occured:!0,details:e.t0};case 9:this.loading=!1;case 10:case"end":return e.stop()}},e,this,[[1,6]])}));function t(){return e.apply(this,arguments)}return t}()},{key:"team",get:function(){return this.teamData(this.teamID)}},{key:"teamID",get:function(){return this.$route.params.team}}]),t}(m["c"]);p["a"]([Object(v["a"])(f["a"].FETCH_TEAM,q)],J.prototype,"fetchTeam",void 0),p["a"]([Object(v["a"])(f["a"].UPDATE_TEAM,q)],J.prototype,"updateTeam",void 0),p["a"]([Object(v["a"])(g["a"].DOWNLOAD_METADATA,$)],J.prototype,"fetchLeagueData",void 0),p["a"]([Object(v["b"])(f["b"].TEAM,{namespace:b["a"].TEAMS})],J.prototype,"teamData",void 0),p["a"]([Object(v["b"])(f["b"].MATCHES,q)],J.prototype,"matchesData",void 0),J=p["a"]([Object(m["a"])({components:{Matches:D,Overview:R}})],J);var K=J,V=K,W=Object(k["a"])(V,n,r,!1,null,null,null);t["default"]=W.exports},"2f21":function(e,t,a){"use strict";var n=a("79e5");e.exports=function(e,t){return!!e&&n(function(){t?e.call(null,function(){},1):e.call(null)})}},"55dd":function(e,t,a){"use strict";var n=a("5ca1"),r=a("d8e8"),i=a("4bf8"),s=a("79e5"),o=[].sort,c=[1,2,3];n(n.P+n.F*(s(function(){c.sort(void 0)})||!s(function(){c.sort(null)})||!a("2f21")(o)),"Array",{sort:function(e){return void 0===e?o.call(i(this)):o.call(i(this),r(e))}})}}]);
//# sourceMappingURL=team.83dd2b71.js.map