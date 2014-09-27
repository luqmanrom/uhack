'use strict';

/**
 * @ngdoc function
 * @name uhackApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the uhackApp
 */
angular.module('uhackApp')
  .controller('MainCtrl', function ($scope,$http) {
        $scope.item = "";
        $scope.stores = ["Target", "Walmart", "Walgreen"];
        $http.defaults.headers.put = {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
            'Access-Control-Allow-Headers': 'Content-Type, X-Requested-With'
        };
        $http.defaults.useXDomain = true;



        $scope.dummies = [ {"stores" : "target","price":200},
                         {"stores" : "walmart","price":200},
                         {"stores" : "walgreen","price":200}
        ];


        console.log($scope.dummies[0].stores);
        console.log($scope.dummies[0].price);


        $scope.submit = function sendData() {
        $scope.show = true;
        $scope.test = $scope.item.split(',').join(' ');
        $scope.url = "http://localhost:8080/rpc?q="+ $scope.test;
        console.log($scope.url);

        $http({
                url: $scope.url,
                method: 'GET'

            })
                .then(function(response) {
                  $scope.respond = angular.fromJson(response).data.totals;
	         $scope.dummies = angular.fromJson(response).data.totals; 
	          
                  console.log(response);
                },
                function(response) {
                    // failed

                }
            );



        };

  });
