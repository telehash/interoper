'use strict';

angular.module('web')
  .directive('logEntry', function () {
    return {
      restrict: 'E',
      scope: {
        entry: "=",
        time: "@",
        sinceTime: "@"
      },
      transclude: true,
      templateUrl: 'components/log-entry/log-entry.html',
      link: function($scope, $elem, $attr) {
        $scope.delta = $scope.time - $scope.sinceTime;
      }
    }
  })
  .filter('timeDelta', function(){
    return function(input){
      var d = input / 1000;
      if (d < 60) return d.toFixed(3) + "s";
      d = d / 60;
      if (d < 60) return d.toFixed(2) + "m";
      d = d / 60;
      return d.toFixed(2) + "h";
    };
  });
