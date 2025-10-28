package networking

import com.google.gson.Gson
import com.google.gson.reflect.TypeToken
import okhttp3.OkHttpClient
import okhttp3.Request
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

@Serializable
data class MmaApiResponse(
    val total: Int,
    val objectIDs: List<Int>
)

data class MmaObjects(val total: Int, val objectIDs: Array<Int>)

@Serializable
data class MmaObject(
    val objectID: Int,
    val isHighlight: Boolean,
    val accessionNumber: String,
    val accessionYear: String,
    val isPublicDomain: Boolean,
    val primaryImage: String,
    val primaryImageSmall: String,
    val additionalImages: List<String>,
    val constituents: List<Constituent>,
    val department: String,
    val objectName: String,
    val title: String,
    val culture: String,
    val period: String,
    val dynasty: String? = null,
    val reign: String? = null,
    val portfolio: String? = null,
    val artistRole: String,
    val artistPrefix: String? = null,
    val artistDisplayName: String,
    val artistDisplayBio: String,
    val artistSuffix: String? = null,
    val artistAlphaSort: String,
    val artistNationality: String,
    val artistBeginDate: String,
    val artistEndDate: String,
    val artistGender: String,
    val artistWikidata_URL: String,
    val artistULAN_URL: String,
    val objectDate: String,
    val objectBeginDate: Int,
    val objectEndDate: Int,
    val medium: String,
    val dimensions: String,
    val measurements: List<Measurement>,
    val creditLine: String,
    val geographyType: String? = null,
    val city: String? = null,
    val state: String? = null,
    val county: String? = null,
    val country: String? = null,
    val region: String? = null,
    val subregion: String? = null,
    val locale: String? = null,
    val locus: String? = null,
    val excavation: String? = null,
    val river: String? = null,
    val classification: String,
    val rightsAndReproduction: String,
    val linkResource: String,
    val metadataDate: String,
    val repository: String,
    val objectURL: String,
    val tags: List<Tag>? = null,
    val objectWikidata_URL: String,
    val isTimelineWork: Boolean,
    val GalleryNumber: String? = null
)

@Serializable
data class Constituent(
    val constituentID: Int,
    val role: String,
    val name: String,
    val constituentULAN_URL: String,
    val constituentWikidata_URL: String,
    val gender: String
)

@Serializable
data class Measurement(
    val elementName: String,
    val elementDescription: String?,
    val elementMeasurements: Map<String, Double>
)

@Serializable
data class Tag(
    val term: String,
    val AAT_URL: String,
    val Wikidata_URL: String
)

@Serializable
data class DepartmentsList(
    val departments: List<Department>,
)

@Serializable
data class Department(
    val departmentId: Int,
    val displayName: String
)

// The Metropolitan Museum of Art Collection API
// https://metmuseum.github.io/
class MmaCollection {
    private val apiUrl = "https://collectionapi.metmuseum.org/public/collection/v1/"
    private val client = OkHttpClient()
    private val gson = Gson()

    private val totalObjects: MmaApiResponse

    init {
        val request = Request.Builder()
            .url("$apiUrl/objects")
            .build()

        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            val respString = response.body?.string()
            throw Exception("bad response: $respString")
        }

        val objectsJson = response.body?.string()

        totalObjects = Json.decodeFromString<MmaApiResponse>(objectsJson!!)
    }

    val objectsCount: Int
        get() {
            return totalObjects.total
        }

    val objectsList: List<Int>
        get() {
            return totalObjects.objectIDs
        }

    fun getObjectAtPosition(id: Int): Int {
        return totalObjects.objectIDs.getOrElse(id) {return -1}
    }

    fun requestObjectData(id: Int) : MmaObject {
        when (id) {
            -1 -> throw Exception("bad object id")
        }

        val request = Request.Builder()
            .url("$apiUrl/objects/$id")
            .build()

        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            val respString = response.body?.string()
            throw Exception("bad response: $respString")
        }

        val objectJson = response.body?.string()

        return Json.decodeFromString<MmaObject>(objectJson!!)
    }

    fun getDepartments(): DepartmentsList {
        val request = Request.Builder()
            .url("$apiUrl/departments")
            .build()

        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            val respString = response.body?.string()
            throw Exception("bad response: $respString")
        }

        val departments = response.body?.string()

        return Json.decodeFromString<DepartmentsList>(departments!!)
    }

    // A listing of all valid Object IDs available for access.
    fun getObjects() {
        println("MmaCollection === get objects list...")

        val request = Request.Builder()
            .url("$apiUrl/objects")
            .build()

        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            val respString = response.body?.string()
            throw Exception("bad response: $respString")
        }

        val objectsJson = response.body?.string()

        val typeToken = object : TypeToken<MmaObjects>() {}.type
        val mmaObjects = gson.fromJson<MmaObjects>(objectsJson, typeToken)

        println(mmaObjects)
    }
}
