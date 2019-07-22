mciModule.factory('RestartUtil', function($filter) {
  return {
    STATUS: {
      ALL: {
        name: 'All',
        matches: _.constant(true),
      },
      NONE: {
        name: 'None',
        matches: _.constant(false),
      },
      FAILURES: {
        name: 'Failures',
        matches: function(task) { return task.status == 'failed' || task.status == 'test-timed-out'}
      },
      SYSTEM_FAILURES: {
        name: 'System Failures',
        matches: function(task) {
          var status = $filter('statusFilter')(task, false);
          return (status == 'system-failed' || status == 'system-unresponsive');
        }
      },
      SETUP_FAILURES: {
        name: 'Setup Failures',
        matches: function(task) {
          var palette = $window.user.Settings.alternate_palette;
          var status = $filter('statusFilter')(task, palette);
          return (status == 'setup-failed' || status == 'setup-failed-alternate')
        }
      },
    }
  }
})
