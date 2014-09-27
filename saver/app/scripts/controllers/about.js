'use strict';

/**
 * @ngdoc function
 * @name uhackApp.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the uhackApp
 */

angular.module('uhackApp')
  .controller('AboutCtrl', function ($scope,myservice) {

//   angular.extend($scope, {
//     center: {
//       autoDiscover: true
//     }
//   });

        $scope.master = {};

        $scope.update = function(user) {
            $scope.master = angular.copy(user);
             myservice.test();
        };

//        $scope.myData = [ { label: "Foo", data: [ [10, 1], [17, -14], [30, 5] ] },
//                          { label: "Bar", data: [ [11, 13], [19, 11], [30, -7] ] }
//                        ];
//        $scope.myChartOptions = {
//            series: {
//                lines: { show: true },
//                points: { show: true }
//            }
//        };





    });
