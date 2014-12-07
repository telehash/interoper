'use strict';

angular.module('web')
  .directive('logLabel', function () {
    return {
      restrict: 'E',
      scope: {
        type: "@",
        identifier: "@",
      },
      templateUrl: 'components/log-label/log-label.html',
      link: function($scope, $elem, $attr) {
        switch ($scope.type) {
        case 'role':
          $scope.formattedType = 'R';
          $scope.formattedIdentifier = ($scope.identifier || "");
          break;
        case 'endpoint':
          $scope.formattedType = 'E';
          $scope.formattedIdentifier = ($scope.identifier || "").slice(0,8);
          break;
        case 'exchange':
          $scope.formattedType = 'X';
          $scope.formattedIdentifier = ($scope.identifier || "").slice(0,8);
          break;
        case 'channel':
          $scope.formattedType = 'C';
          $scope.formattedIdentifier = ($scope.identifier || "");
          break;
        case 'packet':
          $scope.formattedType = 'P';
          $scope.formattedIdentifier = ($scope.identifier || "");
          break;
        }
      }
    }
  });
