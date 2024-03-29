<?php

namespace App\Controllers;


use App\Controllers\BaseController;
use App\Models\AuthModel;

// Clase controladora para el inicio de sesión
class AuthController extends BaseController
{
    private $authModel;

    public function __construct()
    {
        $this->authModel = new AuthModel();
    }

    public function login(): object
    {
        $response = [];
        if (!$this->request->isAJAX()) {
            $response['status'] = 500;
            $response['message'] = 'No se puede procesar la solicitud';
            return $this->response->setJSON($response);
        }

        // Reglas de validación
        $rules = [
            'email' => 'required|valid_email',
            'password' => 'required|min_length[8]|max_length[20]'
        ];

        // Mensajes de validación
        $errors = [
            'email' => [
                'required' => 'El correo electrónico es requerido',
                'valid_email' => 'El correo electrónico no es válido'
            ],
            'password' => [
                'required' => 'La contraseña es requerida',
                'min_length' => 'La contraseña debe tener al menos 8 caracteres',
                'max_length' => 'La contraseña debe tener máximo 20 caracteres'
            ]
        ];

        // Validar datos de entrada
        if (!$this->validate($rules, $errors)) {
            $response['status'] = 500;
            $response['success'] = false;
            $response['message'] = 'No se puede procesar la solicitud';
            $response['errors'] = $this->validator->getErrors();
            return $this->response->setJSON($response);
        }

        $email = $this->request->getPost('email');
        $password = $this->request->getPost('password');

        $user = $this->authModel->where('mail', $email)->first();

        // Verificar si el usuario existe
        if (!$user) {
            $response['status'] = 500;
            $response['success'] = false;
            $response['message'] = 'No se puede procesar la solicitud';
            $response['errors'] = ['email' => 'El correo electrónico no existe'];
            return $this->response->setJSON($response);
        }

        $password_verify = password_verify($password, $user['password']);


        //Validar clave
        if (!$password_verify) {
            $response['status'] = 500;
            $response['success'] = false;
            $response['message'] = 'No se puede procesar la solicitud';
            $response['errors'] = ['password' => 'La contraseña es incorrecta'];
            return $this->response->setJSON($response);
        }

        // Generar token 
        $token = bin2hex(openssl_random_pseudo_bytes(16));

        $session_data = [
            'id' => $user['id'],
            'name' => $user['name'],
            'email' => $user['mail'],
            'isLoggedIn' => true,
            'token' => $token
        ];


        // Retornar data del usuario
        $response['status'] = 200;
        $response['success'] = true;
        $response['message'] = 'Inicio de sesión exitoso';
        $response['data'] = $session_data;
      

        return $this->response->setJSON($response);
    }
}
