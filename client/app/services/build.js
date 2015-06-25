module.exports = function ($http, MessengerService) {
  return {
    getPayloads: getPayloads,
    getPayloadOptions: getPayloadOptions,
    generatePayload: generatePayload,
    getVersion: getVersion
  }

  function requestComplete (response) {
    return response.data
  }

  function getPayloads () {
    return $http.get('/api/payloads')
      .then(requestComplete)
      .catch(MessengerService.error)
  }

  function getPayloadOptions (payload) {
    return $http.get('/api/options?payload=' + payload)
      .then(requestComplete)
      .catch(MessengerService.error)
  }

  function generatePayload (options) {
    return $http.post('/api/generate', {options: options})
      .then(requestComplete)
      .catch(MessengerService.error)
  }

  function getVersion() {
    return $http.get('/api/version')
      .then(requestComplete)
      .catch(MessengerService.error)
  }
}
