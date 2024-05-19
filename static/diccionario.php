<?php
// Conexión a la base de datos
$conexion = mysqli_connect("localhost", "usuario", "contraseña", "nombre_base_de_datos");

// Verificar la conexión
if (mysqli_connect_errno()) {
    echo "Error al conectar a MySQL: " . mysqli_connect_error();
    exit();
}

// Consulta para obtener todas las palabras y descripciones
$consulta = "SELECT palabra, descripcion FROM diccionario WHERE palabra LIKE 'A%'";
$resultado = mysqli_query($conexion, $consulta);

// Crear un array para almacenar los resultados
$palabras = array();
while ($fila = mysqli_fetch_assoc($resultado)) {
    $palabras[] = $fila;
}

// Convertir el array a formato JSON y enviarlo al cliente
echo json_encode($palabras);

// Cerrar la conexión
mysqli_close($conexion);
?>
