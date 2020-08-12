function sleep(milliseconds) {
  var start = new Date().getTime();
  for (var i = 0; i < 1e7; i++) {
    if ((new Date().getTime() - start) > milliseconds){
      break;
    }
  }
}

angular.module('EmailControllers', [])
.controller('PLG', ['$scope', '$http', '$routeParams', function($scope, $http, $routeParams) {
   $scope.schools = getActivities();    // 取得活動資訊
   $scope.getData = 0;
   $scope.sortCondition = ['activity','session'];
   $scope.Headers = ['活動舉辦學校','活動名稱','場次','報名者學校','學生姓名','年級','報名時間'];
   $scope.lists = [];  // 所有資料
   $scope.activities = [];	// 活動列表
   $scope.search = "";	// 搜尋
   $scope.filterCnt = 0;	// filtercnt
   $scope.persons = [];		// 每個人參加的活動
   $scope.showPersonView = 0;

   $scope.Row2String = function(data) {
      var i, d, t, s = $scope.Headers.join(",") + "\n";
      angular.forEach(data, function(d, i) {
         t = [];
         angular.forEach(d, function(val, key) {
            t.push(val);
         });
         s += t.join(",") + "\n";
      });
      return s;
   };

   $scope.downloadExcel = function() {
      var data = $scope.Row2String($scope.lists);
      var form = new FormData();
      form.append('csv', "{}");
      var blob = new Blob(["\ufeff", data], {type:"application/vnd.ms-excel;charset=utf-8"});
      var objectUrl = URL.createObjectURL(blob);
      var anchor = angular.element('<a/>');
      anchor.attr({
         href: objectUrl,
         target: '_blank',
         download: 'youngerboss2020報名表.csv'
      })[0].click();
      form = null;
   };

   // 重設條件
   $scope.ClearFilter = function() {
      $scope.search = "";
      $scope.selectActName = null;
      $scope.filterName = "";
      $scope.filterCnt = 0;
   }; 

   $scope.Backlist = function() { $scope.showPersonView = 0; };

   $scope.setFilter = function() {
      $scope.search = ($scope.selectActName == "" ? "" : $scope.selectActName.name); 
   };

   $scope.sortX = function(s) { 
      if(!s)  $scope.sortCondition = ['activity','session'];
      else $scope.sortCondition = s; 
   };
  
   // 用人名牌活動
   $scope.collectPerson = function() {
      $scope.persons = [];
      var psn = [];
      angular.forEach($scope.lists, function(p) {
         if(psn.indexOf(p["姓名"]) == -1) {
            psn.push(p["姓名"]);
            var x = { name: p["姓名"],school: p["就讀學校"], degree: p["就讀年級"], email: p.Email, send:0, acts:[] };
            x.acts.push({ school: p.school, activity: p.activity, session: p["報名場次"] });
            $scope.persons.push(x); 
         } else {
            for(var i = 0; i < $scope.persons.length && $scope.persons[i].name != p["姓名"]; i++);
            if(i == $scope.persons.length) {
               console.log(p["姓名"] + "找不到");
            } else {
               $scope.persons[i].acts.push({school: p.school, activity: p.activity, session: p["報名場次"] });
            }
         }
      });
      $scope.showPersonView = 1;
   };

   $scope.sendEmail = function() {
      var s = "";
      angular.forEach($scope.persons, function(p) {
       s = "<p>親愛的" + p.name + " 您好</p><p> 恭喜您成功報名8/23 Youngerboss親子技職體驗營課程，敬請於第一堂課半小時前，" +
            "至私立大同高中校門入口報到處，完成報到並領取活動手冊。</p><p>除網路預約課程，現場另有開放體驗課名額，歡迎" +
            "把握資源，勇於嘗試。</p><p>您報名的活動如下：</p>"
          angular.forEach(p.acts, function(a) {
             s += "<li>" + a.session + "&nbsp;" + a.school + "&nbsp;&nbsp; " + a.activity + "</li>\n";
          });
       s += "<p>Ps.體驗課程滿六項，可憑課程章戳，至九樓禮堂換取摸彩券，並於下午4點前至該場地準備參與結業式及摸彩活動。</p>" +
            "<p> 預祝 <br>" + p.name + "探索有成 大獎入手</p><p>臺北市國中學生家長會聯合會<br>總會長 莊孟峯暨全體會員 敬上";

          // "2020 Youngerboss 技職體驗營活動錄取通知信",
          var params = {
            "minetype": "text/html",
            "subject": "=?utf-8?b?MjAyMCBZb3VuZ2VyYm9zcyDmioDogbfpq5TpqZfnh5/mtLvli5XpjITlj5bpgJrnn6Xkv6E=?=",
            "content": s,
            "from": {"Name":"Youngerboss 2020", "email":"andyliu@sinica.edu.tw"},
            "to":  [{"Name": p.name, "email": p.email }],
             // "to":  [{"Name": p.name, "email": "justgps@gmail.com" }],
            "replyto": {"Name": "Andy", "Email":"liuchengood@gmail.com"}
         }
         $http({url:"/email/send", method:"POST", headers: {"Content-Type": "application/json; charset=utf-8"}, data: params })
            .then(function successCallback(res) {
               if(res.data.errMsg) { alert("Data Error: " + res.data.errMsg); return false; }
            console.log(params);  
         }, function errorCallback(data) { console.log(school.name + ": " + act.name + " not found."); return "";});
      }); 
   };

   $scope.init = function() {
      $scope.getData = 1;
      $scope.lists = [];
      var psn, actz, list;
      angular.forEach($scope.schools, function(school, abbr) {
         angular.forEach(school.acts, function(act) {
            list = {school: school.name, activity: act.name};
            actz = { "school": school.name, "name": act.name, timez: [] };
            angular.forEach(act.timez, function(t) {
               actz.timez.push({"time": t});
            });
               var url = "/email/gs?sID=" + school.sheetID + "&sheetName=" + act.name;
               psn = null;
               $http({url: url, method:"GET", headers: {"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"} })
                  .then(function successCallback(res) {
                  if(res.data.errMsg) { alert("Data Error: " + res.data.errMsg); return false; }
                  psn = res.data;
                  angular.forEach(psn.data, function(p) { 
                     p.school = school.name;  p.activity = act.name
                     $scope.lists.push(p);
                  });
                }, function errorCallback(data) { console.log(school.name + ": " + act.name + " not found."); return "";});
            $scope.activities.push(actz);
            list = actz = null;
         });
      });
      $scope.getData = 0;
   };
   $scope.init();
}])
;
