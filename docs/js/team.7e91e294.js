(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["team"],{"0767":function(e,t,a){"use strict";a.r(t);var n=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"home"},[n("img",{attrs:{alt:"Vue logo",src:a("cf05")}}),n("h1",[e._v("Team "+e._s(e.teamID))]),n("p",[e._v("These stats are collected only from games where at least 4 members from this team played together.")]),e.team&&!e.loading?n("div",[n("div",[n("button",{on:{click:function(t){return e.copyMembersToClipboard()}}},[e._v("copy to clipboard")]),n("a",{attrs:{href:"http://na.op.gg/multi/query="+e.memberNames(),target:"_blank"}},[n("button",[e._v("open in na.op.gg")])]),n("button",{on:{click:function(t){return e.forceFetchTeam()}}},[e._v("refresh")])]),n("Overview",{attrs:{teamID:e.teamID}}),n("br"),n("Matches",{attrs:{teamID:e.teamID}}),n("br"),e.updateTriggered?n("div",[e._v("\n      Matches sync queued\n    ")]):e._e(),e.updateTriggered?e._e():n("button",{on:{click:function(t){return e.syncMatches()}}},[e._v("\n      Sync Matches\n    ")])],1):e._e(),e.error.occured?n("div",[e._v("\n    Oops an error occured! "+e._s(e.error.details)+"\n  ")]):e._e(),e.loading?n("div",[e._v("\n    Loading...\n  ")]):e._e()])},r=[],s=(a("7f7f"),a("96cf"),a("3b8d")),i=a("d225"),o=a("b0b4"),c=a("308d"),u=a("6bb5"),l=a("4e2b"),m=a("9ab4"),p=a("60a3"),d=a("4bb5"),v=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"matches"},[a("h2",[e._v("Matches")]),a("p",[e._v("This section is busted right now, don't bother looking.")]),e.matches?e._e():a("div",[e._v("\n    No matches found for this team.\n  ")]),e.matches?a("div",[a("p",[e._v(e._s(e.matches.length)+" matches retrieved.")]),e._l(e.matches,function(t){return a("div",{key:t.gameId},[a("h3",[e._v("Game "+e._s(t.details.gameId))]),a("p",[e._v(e._s(e.dateString(t.details.gameCreation)))]),e._l(t.details.participants,function(t){return a("div",{key:t.participantId},[e._v("\n        Participant: "+e._s(t.participantId)+"\n        "),a("br"),a("div",[a("img",{attrs:{src:e.champIcon(t.championId)}})]),a("br"),a("div",[e._v("\n          Spells:\n          "+e._s(e.spell(t.spell1Id).name)+"\n          "+e._s(e.spell(t.spell2Id).name)+"\n        ")]),a("br"),a("div",[e._v("\n          First Item ("+e._s(t.stats.item0)+"):\n          "),e.item(t.stats.item0)?e._e():a("div",[e._v("\n            No Item\n          ")]),e.item(t.stats.item0)?a("div",[e._v("\n            "+e._s(e.item(t.stats.item0).name)+"\n            "),a("img",{attrs:{src:e.itemIcon(t.stats.item0)}})]):e._e()]),a("br"),a("div",[a("div",[e._v("\n            Primary Perk ("+e._s(t.stats.perkPrimaryStyle)+"):\n            "+e._s(e.runes(t.stats.perkPrimaryStyle).key)+"\n          ")]),a("div",[e._v("\n            Secondary Perk ("+e._s(t.stats.perkSubStyle)+"):\n            "+e._s(e.runes(t.stats.perkSubStyle).key)+"\n          ")])]),a("hr")])})],2)})],2):e._e()])},h=[],b=a("0613"),g=a("8676"),f=a("349a"),y={namespace:b["a"].LEAGUE},_={namespace:b["a"].TEAMS},I=function(e){function t(){return Object(i["a"])(this,t),Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments))}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"dateString",value:function(e){return new Date(e).toLocaleString("en-US")}},{key:"matches",get:function(){return this.matchesData(this.teamID)}}]),t}(p["c"]);m["a"]([Object(p["b"])()],I.prototype,"teamID",void 0),m["a"]([Object(d["b"])(g["b"].MATCHES,_)],I.prototype,"matchesData",void 0),m["a"]([Object(d["b"])(f["b"].ITEM,y)],I.prototype,"item",void 0),m["a"]([Object(d["b"])(f["b"].ITEM_ICON,y)],I.prototype,"itemIcon",void 0),m["a"]([Object(d["b"])(f["b"].CHAMP,y)],I.prototype,"champ",void 0),m["a"]([Object(d["b"])(f["b"].CHAMP_ICON,y)],I.prototype,"champIcon",void 0),m["a"]([Object(d["b"])(f["b"].RUNES,y)],I.prototype,"runes",void 0),m["a"]([Object(d["b"])(f["b"].SPELL,y)],I.prototype,"spell",void 0),I=m["a"]([p["a"]],I);var O=I,j=O,w=a("2877"),k=Object(w["a"])(j,v,h,!1,null,null,null),T=k.exports,D=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"overview"},[a("h2",[e._v("Overview")]),a("div",[a("div",[e._v("\n      "+e._s(e.overviews.team)+"\n    ")]),e._l(e.team.members,function(t){return a("div",{key:t.id},[a("h3",[e._v("\n        "+e._s(t.name)+"\n        ("+e._s(e.overviews.players&&e.overviews.players[t.name]&&e.overviews.players[t.name].tier?e.overviews.players[t.name].tier+", ":"")+"\n        lv"+e._s(t.summonerLevel)+")\n        "),a("a",{attrs:{href:"https://na.op.gg/summoner/userName="+t.name,target:"_blank"}},[a("img",{attrs:{width:"16",height:"16",src:"https://lh3.googleusercontent.com/UdvXlkugn0bJcwiDkqHKG5IElodmv-oL4kHlNAklSA2sdlVWhojsZKaPE-qFPueiZg"}})])]),e.overviews.players&&e.overviews.players[t.name]&&e.overviews.players[t.name].aggs?a("div",[a("div",[a("h5",[e._v("Most played lane and role")]),e._v("\n          "+e._s(e.overviews.players[t.name].aggs.favourite.lane)+"\n          ("+e._s(e.overviews.players[t.name].aggs.favourite.role)+")\n        ")]),a("div",[a("h5",[e._v("Most played champions")]),e._l(e.overviews.players[t.name].aggs.favourite.champs,function(n){return a("img",{key:"fav-"+t.name+"-"+n,attrs:{src:e.champIcon(n)}})})],2),a("div",[a("h5",[e._v("Average stats")]),e._v("\n          "+e._s(e.overviews.players[t.name].aggs.avg)+"\n        ")])]):e._e()])})],2)])},E=[];a("456d"),a("7514"),a("ac6a"),a("55dd");function M(){return{vision:[],cs:[],jungle:{friendly:[],enemy:[]},dealt:[],taken:[],gold:[],champs:[],lanes:[],roles:[]}}function A(e,t){e.vision.push(t.stats.visionScore),e.cs.push(t.stats.totalMinionsKilled),e.cs.push(t.stats.neutralMinionsKilled),e.jungle.friendly.push(t.stats.neutralMinionsKilledTeamJungle),e.jungle.enemy.push(t.stats.neutralMinionsKilledEnemyJungle),e.dealt.push(t.stats.totalDamageDealtToChampions),e.taken.push(t.stats.totalDamageTaken),e.gold.push(t.stats.goldEarned),e.champs.push(t.championId),e.lanes.push(t.timeline.lane),e.roles.push(t.timeline.role)}function S(e){var t=e.reduce(function(e,t){return e+t},0);return(t/e.length).toFixed(2)}function C(e,t){var a=e.reduce(function(e,t){return e[t]=(e[t]||0)+1,e},{}),n=Object.keys(a).sort(function(e,t){return a[t]-a[e]});return t?n.slice(0,t):n[0]}function N(e){return{avg:{vision:S(e.vision),cs:S(e.cs),dealt:S(e.dealt),taken:S(e.taken),gold:S(e.gold),jungle:{friendly:S(e.jungle.friendly),enemy:S(e.jungle.enemy)}},favourite:{champs:C(e.champs,5),lane:C(e.lanes),role:C(e.roles)}}}function P(e){var t=e.wins.filter(function(e){return e}).length,a=parseInt(S(e.matchTime),10);return{games:e.wins.length,winRate:(t/e.wins.length).toFixed(2),avg:{matchTime:"".concat(Math.floor(a/60),":").concat(Math.floor(a%60)),towers:S(e.towers),dragons:S(e.dragons),barons:S(e.barons)}}}var x={namespace:b["a"].LEAGUE},L={namespace:b["a"].TEAMS},H=function(e){function t(){return Object(i["a"])(this,t),Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments))}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"idToName",value:function(){var e={};return this.teamData(this.teamID).members.forEach(function(t){e[t.accountId]=t.name}),e}},{key:"team",get:function(){return this.teamData(this.teamID)}},{key:"overviews",get:function(){var e=this.matchesData(this.teamID);if(!e)return{team:{},players:{}};var t=this.idToName();console.debug("generating overviews",{matches:e.length});var a={},n={matchTime:[],wins:[],towers:[],dragons:[],barons:[]},r={team:{},players:{}};return e.forEach(function(e){var s;e.details.participantIdentities.forEach(function(n){var i=n.player.accountId;n.player.currentAccountId&&n.player.currentAccountId!==i&&(i=n.player.currentAccountId);var o=t[i];if(o){var c=e.details.participants.find(function(e){return e.participantId===n.participantId});if(c){var u=e.details.teams.find(function(e){return e.teamId===c.teamId});u&&(s=u),r.players[o]||(r.players[o]={}),r.players[o].tier=c.highestAchievedSeasonTier,a[o]||(a[o]=M()),A(a[o],c)}else console.debug("could not find participant ".concat(n.participantId))}}),n.matchTime.push(e.details.gameDuration),s&&(n.wins.push("Win"===s.win),n.towers.push(s.towerKills),n.dragons.push(s.dragonKills),n.barons.push(s.baronKills))}),Object.keys(a).forEach(function(e){r.players[e].aggs=N(a[e])}),r.team.aggs=P(n),console.log("generated overviews",r),r}}]),t}(p["c"]);m["a"]([Object(p["b"])()],H.prototype,"teamID",void 0),m["a"]([Object(d["b"])(g["b"].MATCHES,L)],H.prototype,"matchesData",void 0),m["a"]([Object(d["b"])(g["b"].TEAM,{namespace:b["a"].TEAMS})],H.prototype,"teamData",void 0),m["a"]([Object(d["b"])(f["b"].ITEM,x)],H.prototype,"item",void 0),m["a"]([Object(d["b"])(f["b"].ITEM_ICON,x)],H.prototype,"itemIcon",void 0),m["a"]([Object(d["b"])(f["b"].CHAMP,x)],H.prototype,"champ",void 0),m["a"]([Object(d["b"])(f["b"].CHAMP_ICON,x)],H.prototype,"champIcon",void 0),m["a"]([Object(d["b"])(f["b"].RUNES,x)],H.prototype,"runes",void 0),m["a"]([Object(d["b"])(f["b"].SPELL,x)],H.prototype,"spell",void 0),H=m["a"]([p["a"]],H);var K=H,F=K,U=Object(w["a"])(F,D,E,!1,null,null,null),R=U.exports;function G(e){var t=document.createElement("textarea");t.value=e,t.setAttribute("readonly",""),document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t)}var J={namespace:b["a"].TEAMS},q={namespace:b["a"].LEAGUE},$=function(e){function t(){var e;return Object(i["a"])(this,t),e=Object(c["a"])(this,Object(u["a"])(t).apply(this,arguments)),e.error={occured:!1},e.loading=!0,e.updateTriggered=!1,e}return Object(l["a"])(t,e),Object(o["a"])(t,[{key:"mounted",value:function(){var e=Object(s["a"])(regeneratorRuntime.mark(function e(){return regeneratorRuntime.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,this.fetchTeam({teamID:this.teamID});case 3:return e.next=5,this.fetchLeagueData({});case 5:e.next=10;break;case 7:e.prev=7,e.t0=e["catch"](0),this.error={occured:!0,details:e.t0};case 10:this.loading=!1;case 11:case"end":return e.stop()}},e,this,[[0,7]])}));function t(){return e.apply(this,arguments)}return t}()},{key:"syncMatches",value:function(){this.updateTriggered=!0,this.error={occured:!1};try{this.updateTeam({teamID:this.teamID})}catch(e){this.error={occured:!0,details:e},this.updateTriggered=!1}}},{key:"memberNames",value:function(){return this.team?this.team.members.map(function(e){return e.name}).join(","):""}},{key:"copyMembersToClipboard",value:function(){var e=this.memberNames();G(e)}},{key:"forceFetchTeam",value:function(){var e=Object(s["a"])(regeneratorRuntime.mark(function e(){return regeneratorRuntime.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return this.loading=!0,e.prev=1,e.next=4,this.fetchTeam({teamID:this.teamID,force:!0});case 4:e.next=9;break;case 6:e.prev=6,e.t0=e["catch"](1),this.error={occured:!0,details:e.t0};case 9:this.loading=!1;case 10:case"end":return e.stop()}},e,this,[[1,6]])}));function t(){return e.apply(this,arguments)}return t}()},{key:"team",get:function(){return this.teamData(this.teamID)}},{key:"teamID",get:function(){return this.$route.params.team}}]),t}(p["c"]);m["a"]([Object(d["a"])(g["a"].FETCH_TEAM,J)],$.prototype,"fetchTeam",void 0),m["a"]([Object(d["a"])(g["a"].UPDATE_TEAM,J)],$.prototype,"updateTeam",void 0),m["a"]([Object(d["a"])(f["a"].DOWNLOAD_METADATA,q)],$.prototype,"fetchLeagueData",void 0),m["a"]([Object(d["b"])(g["b"].TEAM,{namespace:b["a"].TEAMS})],$.prototype,"teamData",void 0),m["a"]([Object(d["b"])(g["b"].MATCHES,J)],$.prototype,"matchesData",void 0),$=m["a"]([Object(p["a"])({components:{Matches:T,Overview:R}})],$);var W=$,V=W,Z=Object(w["a"])(V,n,r,!1,null,null,null);t["default"]=Z.exports},"2f21":function(e,t,a){"use strict";var n=a("79e5");e.exports=function(e,t){return!!e&&n(function(){t?e.call(null,function(){},1):e.call(null)})}},"55dd":function(e,t,a){"use strict";var n=a("5ca1"),r=a("d8e8"),s=a("4bf8"),i=a("79e5"),o=[].sort,c=[1,2,3];n(n.P+n.F*(i(function(){c.sort(void 0)})||!i(function(){c.sort(null)})||!a("2f21")(o)),"Array",{sort:function(e){return void 0===e?o.call(s(this)):o.call(s(this),r(e))}})}}]);
//# sourceMappingURL=team.7e91e294.js.map