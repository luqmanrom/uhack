'use strict';

/**
 * @ngdoc overview
 * @name uhackApp
 * @description
 * # uhackApp
 *
 * Main module of the application.
 */
angular
  .module('uhackApp', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });







