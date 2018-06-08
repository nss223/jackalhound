package org.lab2528.multi

import java.util.Collection

import org.hyperledger.fabric.sdk.ProposalResponse

object Run extends App {
    Util.client.setUserContext(UserManagement.getOrCreateUser("hfuser"))
    val res: Collection[ProposalResponse] = Util.query("multi", "get", Array("a"))

    import scala.collection.JavaConversions._

    for (pres <- res) {
      val stringResponse: String = new String(pres.getChaincodeActionResponsePayload)
      Util.log.info(stringResponse)
    }
}