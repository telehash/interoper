'use strict';

angular.module('web')
  .controller('TestShowCtrl', function ($scope, test) {
    $scope.test = test;
  });
