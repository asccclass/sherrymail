angular.module('Assnx', [
   "ngRoute",
   "EmailControllers"
])
.config(function($routeProvider) {
   $routeProvider
   .when('/', { templateUrl: 'views/index.html'})
   .when('/all', { controller: 'PLG', templateUrl: 'views/all.html'})      // 祕書處
   .otherwise({ redirectTo: '/' });
});
