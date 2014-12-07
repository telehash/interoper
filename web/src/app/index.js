'use strict';

angular.module('web', ['ngAnimate', 'ngCookies', 'ngTouch', 'ngSanitize', 'ngResource', 'ui.router'])
  .config(function ($stateProvider, $urlRouterProvider, $locationProvider) {
    $stateProvider

      .state('tests', {
        url: '/',
        templateUrl: 'app/test/list.html',
        controller: 'TestListCtrl',
        resolve: {
          dump: function($http){
            return $http.get('/dump.json').then(function(resp){ return resp.data; });
          }
        }
      })

      .state('tests.test', {
        url: 'tests/:testName',
        templateUrl: 'app/test/show.html',
        controller: 'TestShowCtrl',
        resolve: {
          test: function(dump, $stateParams){
            return dump.tests[$stateParams.testName];
          }
        }
      })

      .state('tests.test.run', {
        url: '/runs/:driverName',
        templateUrl: 'app/inspector/run.html',
        controller: 'InspectorRunCtrl',
        resolve: {
          run: function(test, $stateParams){
            return test.runs[$stateParams.driverName];
          }
        }
      });

    $urlRouterProvider.otherwise('/');
    $locationProvider.html5Mode(true);
  })
;
