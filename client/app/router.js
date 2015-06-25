module.exports = function ($routeProvider) {
  $routeProvider
    .when('/build', {
      templateUrl: 'views/build.html',
      controller: 'BuildController',
      controllerAs: 'vm'
    })
    .otherwise({
      redirectTo: '/build'
    })
}
