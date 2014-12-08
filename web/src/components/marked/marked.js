'use strict';

angular.module('web')
  .directive('markedElement', function () {
    return {
      restrict: 'E',
      scope: {
        text: '='
      },
      link: function($scope, $elem) {
        $elem.html(marked($scope.text));

        $scope.$watch('text', function(v){
          $elem.html(marked(v));
        });
      }
    };
  });
