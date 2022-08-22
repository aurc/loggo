// package: loggo
// file: loggo/service.proto

var loggo_service_pb = require("../loggo/service_pb");
var loggo_types_pb = require("../loggo/types_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Loggo = (function () {
  function Loggo() {}
  Loggo.serviceName = "loggo.Loggo";
  return Loggo;
}());

Loggo.LogStream = {
  methodName: "LogStream",
  service: Loggo,
  requestStream: false,
  responseStream: true,
  requestType: loggo_types_pb.LogEntryRequest,
  responseType: loggo_types_pb.LogEntryResponse
};

exports.Loggo = Loggo;

function LoggoClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

LoggoClient.prototype.logStream = function logStream(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(Loggo.LogStream, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

exports.LoggoClient = LoggoClient;

