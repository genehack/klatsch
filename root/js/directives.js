var klatschApp = angular.module('klatsch',['ui.bootstrap'],function() {});

klatschApp.directive("countedMax", function() {
  return  {
    require: "ngModel" ,
    restrict: "A" ,
    link: function(scope,element,attrs,ngModelCtrl) {
      var maxlength = Number(attrs.countedMax);
      function fromUser(text) {
        ngModelCtrl.$setValidity('unique', text.length <= maxlength);
        return text;
      }
      ngModelCtrl.$parsers.unshift(fromUser);
    }
  };
}).directive("twitterInput", function () {
  return {
    restrict: "E" ,
    replace: true ,
    templateUrl: "./twitter_input.html" ,
    controller: function($scope,$timeout) {
      $scope.error = function(name) {
        var form = $scope.form[name];
        return form.$invalid && form.$dirty ? "has-error" : "";
      };
    },
    scope: {}
  };
});
