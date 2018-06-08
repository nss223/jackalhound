package org.lab223.multi;

import org.lab2528.multi.UserManagement;
import org.lab2528.multi.Util;

public class Newchannelname {
    public static void main(String[] args) {
        try {

            // invoke & query demo

            StringBuilder txid;
            String res;

            Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));

            txid = new StringBuilder();
            if (Util.invoke("mychannel", "multicc", "set",
                    new String[] { "A", "A", "300" }, txid))
                Util.log.info(txid);

            res = Util.query("mychannel", "multicc", "get", new String[] { "A", "A" });
            Util.log.info(res);
        } catch (Exception e) {
            Util.log.error(e.getMessage());
        }
    }
}
