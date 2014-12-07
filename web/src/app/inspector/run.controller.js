'use strict';

angular.module('web')
  .controller('InspectorRunCtrl', function ($scope, $controller, run) {
    analyzeLog(run);

    $scope.sinceTime = run && run.log && run.log.length && run.log[0].ti;
    $scope.run = run;

    $scope.templateForEntry = function templateForEntry(entry) {
      switch (entry.ty) {
      case 'ready':
        return 'app/inspector/entry-ready.html';
      case 'done':
        return 'app/inspector/entry-done.html';
      case 'exited':
        return 'app/inspector/entry-exited.html';
      case 'log':
        return 'app/inspector/entry-log.html';
      case 'exec':
        return 'app/inspector/entry-exec.html';
      case 'status':
        return 'app/inspector/entry-status.html';

      case 'endpoint.new':
        return 'app/inspector/entry-endpoint-new.html';
      case 'endpoint.started':
        return 'app/inspector/entry-endpoint-started.html';
      case 'endpoint.rcv.packet':
        return 'app/inspector/entry-endpoint-rcv-packet.html';

      case 'exchange.new':
        return 'app/inspector/entry-exchange-new.html';
      case 'exchange.started':
        return 'app/inspector/entry-exchange-started.html';
      case 'exchange.stopped':
        return 'app/inspector/entry-exchange-stopped.html';
      case 'exchange.rcv.handshake':
        return 'app/inspector/entry-exchange-rcv-handshake.html';
      case 'exchange.rcv.packet':
        return 'app/inspector/entry-exchange-rcv-packet.html';


      case 'channel.new':
        return 'app/inspector/entry-channel-new.html';
      case 'channel.write':
        return 'app/inspector/entry-channel-write.html';
      case 'channel.rcv.packet':
        return 'app/inspector/entry-channel-rcv-packet.html';

      default:
        return 'app/inspector/entry-unknown.html';
      }
    };
  })
  .filter('ansi', function($sce){
    return function(input) {
      input = ansi_up.escape_for_html(input);
      input = ansi_up.ansi_to_html(input, {use_classes: true});
      return $sce.trustAsHtml(input);
    };
  })
  .filter('linesplit', function(){
    return function(input) {
      var lines = [];
      while (input.length > 80) {
        lines.push(input.slice(0, 80));
        input = input.slice(80);
      }
      if (input.length > 0) {
        lines.push(input);
      }
      return lines.join("\n");
    };
  });

function analyzeLog(run) {

  run.$$cache = { roles: {} };

  run.log.forEach(function(entry){
    if (entry.in) {
      if (entry.in.endpoint_id) {
        entry.$$endpoint = getEndpoint(entry.ro, entry.in.endpoint_id);
      }

      if (entry.in.exchange_id) {
        entry.$$exchange = getExchange(entry.ro, entry.in.exchange_id);
      }

      if (entry.in.channel_id) {
        entry.$$channel = getChannel(entry.ro, entry.in.channel_id);
      }

      if (entry.in.packet_id) {
        entry.$$packet = getPacket(entry.ro, entry.in.packet_id);
      }
    }

    switch (entry.ty) {

    case 'endpoint.new':
      entry.$$endpoint.hashname = entry.in.hashname;
      break;
    case 'endpoint.started':
      break;
    case 'endpoint.rcv.packet':
      entry.$$packet.endpoint = entry.$$endpoint;
      entry.$$packet.outer    = entry.in.packet;
      break;

    case 'exchange.new':
      entry.$$endpoint.exchanges[entry.in.exchange_id] = entry.$$exchange;
      entry.$$exchange.endpoint = entry.$$endpoint;
      break;
    case 'exchange.started':
      entry.$$exchange.peer = entry.in.peer;
      break;
    case 'exchange.rcv.packet':
      entry.$$packet.exchange = entry.$$exchange;
      entry.$$packet.inner    = entry.in.packet;
    case 'exchange.rcv.handshake':
      entry.$$packet.exchange  = entry.$$exchange;
      entry.$$packet.handshake = entry.in.handshake;
      break;

    case 'channel.new':
      entry.$$exchange.channels[entry.in.channel_id] = entry.$$channel;
      entry.$$exchange.endpoint.channels[entry.in.channel_id] = entry.$$channel;
      entry.$$channel.exchange = entry.$$exchange;
      entry.$$channel.endpoint = entry.$$exchange.endpoint;
      break;
    case 'channel.write':
      entry.$$packet.channel  = entry.$$channel;
      entry.$$packet.exchange = entry.$$channel.exchange;
      entry.$$packet.endpoint = entry.$$channel.endpoint;
      entry.$$packet.inner    = entry.in.packet;
      break;
    case 'channel.rcv.packet':
      entry.$$packet.channel  = entry.$$channel;
      break;

    }
  });

  run.log.forEach(function(entry){
    if (!entry.$$channel && entry.$$packet) {
      entry.$$channel = entry.$$packet.channel;
    }
    if (!entry.$$exchange && entry.$$packet) {
      entry.$$exchange = entry.$$packet.exchange;
    }
    if (!entry.$$endpoint && entry.$$packet) {
      entry.$$endpoint = entry.$$packet.endpoint;
    }

    if (!entry.$$endpoint && entry.$$exchange) {
      entry.$$endpoint = entry.$$exchange.endpoint;
    }

  });

  function getPacket(role, id) {
    role = getRole(role);

    var e = role.packets[id];
    if (!e) {
      e = {
        id: id,
      }
      role.packets[id] = e;
    }
    return e;
  }

  function getChannel(role, id) {
    role = getRole(role);

    var e = role.channels[id];
    if (!e) {
      e = {
        id: id,
      }
      role.channels[id] = e;
    }
    return e;
  }

  function getExchange(role, id) {
    role = getRole(role);

    var e = role.exchanges[id];
    if (!e) {
      e = {
        id: id,
        channels:  {}
      }
      role.exchanges[id] = e;
    }
    return e;
  }

  function getEndpoint(role, id) {
    role = getRole(role);

    var e = role.endpoints[id];
    if (!e) {
      e = {
        id:        id,
        exchanges: {},
        channels:  {}
      }
      role.endpoints[id] = e;
    }
    return e;
  }

  function getRole(role) {
    var r = run.$$cache.roles[role];
    if (!r) {
      r = {
        name:      role,
        endpoints: {},
        exchanges: {},
        channels:  {},
        packets:   {}
      };
      run.$$cache.roles[role] = r;
    }
    return r;
  }
}
