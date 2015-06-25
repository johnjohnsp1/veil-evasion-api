var $ = require('jquery')

module.exports = function (BuildService, MessengerService) {
  var vm = this
  vm.payload = ''
  vm.payloads = []
  vm.payloadOptions = []
  vm.payloadSelected = false
  vm.payloadGenerated = false
  vm.generatedPayload = ''
  vm.getPayloadOptions = getPayloadOptions
  vm.generatePayload = generatePayload
  vm.reset = reset

  activate()

  function activate () {
    BuildService.getPayloads()
      .then(function (data) {
        vm.payloads = data
      })
  }

  function reset() {
    vm.payload = ''
    vm.payloadOptions = []
    vm.payloadSelected = false
    vm.payloadGenerated = false
    vm.generatedPayload = ''
  }

  function getPayloadOptions () {
    BuildService.getPayloadOptions(vm.payload)
      .then(function (data) {
        vm.payloadOptions = data
        vm.payloadSelected = true
      })
  }

  function generatePayload() {
    var options = []
    options.push({
      key: 'payload',
      value: vm.payload
    })
    $('form#generate-form :input').each(function () {
      var val = $(this).val()
      var key = this.name
      if (typeof key !== 'string' || key === '') {
        return
      }
      if (typeof val !== 'string' || val === '') {
        return
      }
      options.push({
        key: key,
        value: val
      })
    })
    BuildService.generatePayload(options)
      .then(function (data) {
        if (typeof data.result === 'undefined') {
          return
        }
        vm.generatedPayload = data.result
        vm.payloadGenerated = true
      })
  }

}
