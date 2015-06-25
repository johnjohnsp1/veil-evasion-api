module.exports = function ($rootScope, BuildService) {
  var vm = this
  vm.dismiss = dismiss
  vm.version = ''

  activate()

  function activate() {
    BuildService.getVersion()
      .then(function (data) {
        vm.version = data.version
      })
  }

  function dismiss () {
    $rootScope.errorMessage = ''
    $rootScope.successMessage = ''
  }
}
