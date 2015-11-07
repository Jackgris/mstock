angular
    .module('mstock')
    .controller('HomeController', HomeController);

function HomeController($log){
    $log.info('Test');
}
