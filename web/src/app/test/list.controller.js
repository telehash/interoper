'use strict';

angular.module('web')
  .controller('TestListCtrl', function ($scope, dump) {
    $scope.dump = dump;
  });
