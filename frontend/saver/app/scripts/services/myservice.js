angular.module('uhackApp')
    .factory('myservice', function($rootScope) {

    var test = function() {
        console.log("success");
    };




     return {
        test: test
     };


});