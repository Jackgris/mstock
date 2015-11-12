angular
    .module('mstock')
    .controller('HomeController', HomeController)
    .controller('LoginController', LoginController)
    .controller('LogoutController', LogoutController)
    .controller('SignUpController', SignUpController);

function HomeController($log){
    $log.info('Test');
}

function SignUpController($auth, $location, $scope, $log) {  

    $scope.signup = function() {
        $auth.signup({
            email: $scope.signup.email,
            password: $scope.signup.password,
            name: $scope.signup.nombre
        })
        .then(function() {
            $log.info('Se realizo el registro correctamente');
            // Si se ha registrado correctamente,
            // Podemos redirigirle a otra parte
            $location.path("/home");
        })
        .catch(function(response) {
            // Si ha habido errores, llegaremos a esta función
            $log.info('Hubo un error en el registro');
        });
    }
}

function LoginController($log, $auth, $location, $scope) {  
    $log.debug('Se carga el controller de login');
    $scope.login = function(){
        $auth.login({
            email: $scope.login.email,
            password: $scope.login.password
        })
        .then(function(){
            // Si se ha logueado correctamente, lo tratamos aquí.
            // Podemos también redirigirle a una ruta
            $log.info('Se realizo el login correctamente');
            $location.path("/home")
        })
        .catch(function(response){
            // Si ha habido errores llegamos a esta parte
            $log.info('Hubo un error en el login');
            $location.path("/registro")
        });
    }
}

function LogoutController($auth, $location) {  
    $auth.logout()
        .then(function() {
            // Desconectamos al usuario y lo redirijimos
            $location.path("/")
        });
}
