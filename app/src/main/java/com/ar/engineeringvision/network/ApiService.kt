package com.ar.engineeringvision.network

import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Path

// Соответствует структуре Element в Go-бэкенде
data class ElementDto(
    val id: String,
    val type: String, // "duct", "pipe", "socket"
    val world_coords: Map<String, Double>, // {"x": 1.2, "y": 0.9, "z": 0.0}
    val properties: Map<String, Any>
)

interface ApiService {
    @GET("api/v1/rooms/{room_id}/elements")
    suspend fun getRoomElements(
        @Header("Authorization") token: String,
        @Path("room_id") roomId: String
    ): List<ElementDto>
}
