package networking

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.main
import com.github.ajalt.clikt.parameters.options.help
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.prompt

class Networking : CliktCommand() {
    // private val userToken: String by option("-t").prompt().help("security token")

    override fun run() {
        try {
            val mma = MmaCollection()
            println(mma.objectsCount)
            // println(mma.objectsList)

            val objectInfo = mma.requestObjectData(mma.getObjectAtPosition(100))
            println(objectInfo)

            println(mma.getDepartments())
        } catch (e: Exception) {
            println("exception was: ${e.message}")
        }
    }
}

fun main(args: Array<String>) {
    Networking().main(args)
}

