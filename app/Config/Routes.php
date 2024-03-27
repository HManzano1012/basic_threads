<?php

use CodeIgniter\Router\RouteCollection;
use App\Controllers\AuthController;

/**
 * @var RouteCollection $routes
 */



$routes->get('/', 'Home::index');
$routes->post('/login', 'AuthController::login');
