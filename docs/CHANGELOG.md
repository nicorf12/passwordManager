# Historial de Versiones

## 0.3.0: Configuración y sincronización

**Versión actualizada del programa.**  
Se centra en la sincronización de datos entre distintas instancias del programa y configuración:

- **Archivos de configuración** para ajustes como idioma, encriptado por defecto, y preferencias de interfaz como tema.
- **Sincronización de contraseñas** entre dos computadoras no vinculadas mediante importar/exportar archivos.

El archivo se encripta 3 veces con 3 claves derivadas de respuestas a preguntas genéricas a la hora de importar/exportar.

---

## 0.2.0 : Mejora
**Versión actualizada del programa.**  
Se cebtra en agregar nuevas funcionalidades y complementar con un nuevo diseño de interfaz:
- **Tipos de encriptado:** Compatibilidad con distintos métodos de encriptado.
- **Información adicional:** Más detalles asociados a cada contraseña almacenada.
- **Nivel de seguridad:** Indicador del nivel de seguridad de las contraseñas.
- **Carpetas personalizadas:** Organización manual de contraseñas en carpetas.
- **Sección de favoritos:** Acceso rápido a contraseñas marcadas como favoritas.
- **Funcionamiento en segundo plano:** El programa puede ejecutarse en segundo plano para mayor comodidad.

---

## 0.1.0 : MVP
Versión mínima viable con funciones básicas para un gestor de contraseñas:
- **Gestión de cuentas:** Registro y inicio de sesión.
- **Contraseñas seguras:** Registro, edición, copia y generación automática de contraseñas seguras.
- **Encriptado y seguridad:**
    - Encriptado AES para proteger contraseñas.
    - Hashing PBKDF2 para proteger contraseñas maestras.
- **Etiquetas:** Sistema básico de etiquetas para identificar contraseñas.
- **Persistencia y multilenguaje:** Mantiene la sesión activa y ofrece soporte para múltiples idiomas.  
