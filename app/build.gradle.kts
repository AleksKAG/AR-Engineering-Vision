dependencies {
    // ARCore
    implementation("com.google.ar:core:1.40.0")
    
    // Sceneform (поддержка glTF/glb)
    implementation("com.gorisse.thomas.sceneform:sceneform:1.21.0")
    
    // Сетевой клиент для получения данных с Go-бэкенда
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    
    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")
}
