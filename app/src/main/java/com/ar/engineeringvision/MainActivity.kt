package com.ar.engineeringvision

import android.Manifest
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Bundle
import android.view.MotionEvent
import android.view.View
import android.widget.Button
import android.widget.TextView
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.google.ar.core.Anchor
import com.google.ar.core.HitResult
import com.google.ar.core.Plane
import com.google.ar.sceneform.rendering.ModelRenderable
import com.google.ar.sceneform.ux.ArFragment
import com.google.ar.sceneform.ux.BaseArFragment

class MainActivity : AppCompatActivity() {

    private lateinit var arFragment: ArFragment
    private lateinit var infoTextView: TextView
    private var currentAnchor: Anchor? = null

    // Заглушка токена (в реальном приложении приходит после Login на Go-бэкенде)
    private val authToken = "Bearer YOUR_JWT_TOKEN_HERE"
    private val roomId = "550e8400-e29b-41d4-a716-446655440000"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        infoTextView = findViewById(R.id.tv_info_overlay)
        arFragment = supportFragmentManager.findFragmentById(R.id.ar_fragment) as ArFragment

        setupButtons()
        checkCameraPermission()
    }

    private fun setupButtons() {
        findViewById<Button>(R.id.btn_photo).setOnClickListener {
            Toast.makeText(this, "Фото сохранено и отправлено на /api/v1/elements/{id}/issues", Toast.LENGTH_SHORT).show()
        }
        findViewById<Button>(R.id.btn_check).setOnClickListener {
            Toast.makeText(this, "Запуск AI-проверки смещения...", Toast.LENGTH_SHORT).show()
        }
        findViewById<Button>(R.id.btn_layers).setOnClickListener {
            infoTextView.text = "Слой: ОВ (Вентиляция)\nСлой: ЭОМ (Электрика)"
        }
        findViewById<Button>(R.id.btn_back).setOnClickListener {
            finish()
        }
    }

    private fun checkCameraPermission() {
        if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) != PackageManager.PERMISSION_GRANTED) {
            ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.CAMERA), 100)
        } else {
            setupArFragment()
        }
    }

    override fun onRequestPermissionsResult(requestCode: Int, permissions: Array<out String>, grantResults: IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == 100 && grantResults.isNotEmpty() && grantResults[0] == PackageManager.PERMISSION_GRANTED) {
            setupArFragment()
        } else {
            Toast.makeText(this, "Требуется доступ к камере для работы AR", Toast.LENGTH_LONG).show()
            finish()
        }
    }

    private fun setupArFragment() {
        // Отключаем стандартные подсказки Sceneform для чистого UI
        arFragment.planeDiscoveryController.hide()
        arFragment.planeDiscoveryController.setInstructionView(null)

        arFragment.setOnTapArPlaneListener { hitResult: HitResult, plane: Plane, motionEvent: MotionEvent ->
            // Реагируем только на горизонтальные плоскости (пол/потолок)
            if (plane.type != Plane.Type.HORIZONTAL_UPWARD_FACING && plane.type != Plane.Type.HORIZONTAL_DOWNWARD_FACING) {
                return@setOnTapArPlaneListener
            }

            // 1. Создаем Anchor (точку привязки в реальном мире)
            val anchor = hitResult.createAnchor()
            currentAnchor = anchor

            // 2. URL модели. В реальности это Presigned URL из MinIO, полученный от Go-бэкенда.
            // Для теста используем открытую glTF модель (утка или простой куб)
            val modelUrl = "https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf"
            
            // 3. Загружаем и рендерим модель
            ModelRenderable.builder()
                .setSource(this, Uri.parse(modelUrl))
                .build()
                .thenAccept { renderable ->
                    val node = com.google.ar.sceneform.Node().apply {
                        this.anchor = anchor
                        this.renderable = renderable
                        // Масштабируем для наглядности (утка большая по умолчанию)
                        localScale = com.google.ar.sceneform.math.Vector3(0.05f, 0.05f, 0.05f)
                    }
                    arFragment.arSceneView.scene.addChild(node)

                    // 4. Обновляем UI согласно ТЗ
                    updateOverlayUI("Воздуховод 600x300", "72 мм")
                }
                .exceptionally { throwable ->
                    Toast.makeText(this, "Ошибка загрузки модели: ${throwable.message}", Toast.LENGTH_LONG).show()
                    null
                }
        }
    }

    private fun updateOverlayUI(elementName: String, deviation: String) {
        runOnUiThread {
            infoTextView.text = """
                --------------------------
                | $elementName           |
                | Кабельный лоток 100 мм |
                | Розетка 900 мм         |
                --------------------------
                ⚠️ Ошибка: смещение $deviation
            """.trimIndent()
        }
    }

    // TODO: Шаг 3 интеграции - раскомментировать и использовать для получения реальных координат
    /*
    private fun loadElementsFromBackend() {
        lifecycleScope.launch {
            try {
                val retrofit = Retrofit.Builder()
                    .baseUrl("http://10.0.2.2:8080/") // 10.0.2.2 - это localhost для Android Emulator
                    .addConverterFactory(GsonConverterFactory.create())
                    .build()
                val api = retrofit.create(ApiService::class.java)
                
                val elements = api.getRoomElements(authToken, roomId)
                for (element in elements) {
                    // Здесь логика размещения Node по element.world_coords
                }
            } catch (e: Exception) {
                Toast.makeText(this@MainActivity, "Ошибка сети: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }
    */
}
