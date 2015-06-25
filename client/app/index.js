var angular = require('angular')
require('./styles/main.styl')
require('bootstrap-webpack')

var router = require('./router')
var controllers = require('./controllers')
var services = require('./services')

var app = angular.module('app', [
  require('angular-sanitize'),
  require('angular-route'),
  require('angular-animate')
])

app.factory('MessengerService', ['$rootScope', '$timeout', '$location', services.Messenger])
app.factory('BuildService', ['$http', 'MessengerService', services.Build])

app.controller('ApplicationController', ['$rootScope', 'BuildService', controllers.App])
app.controller('BuildController', ['BuildService', 'MessengerService', controllers.Build])

app.config(['$routeProvider', router])
