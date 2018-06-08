package org.lab223.multi;

import org.lab2528.multi.UserManagement;
import org.lab2528.multi.Util;

public class TestSplit {
    public static void main(String[] args) {

        // Demo 3
        // split asset into smaller ones

        StringBuilder txid;
        String res;

        try {
            // set admin context
            Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));

            Util.log.info("Register user `user_sp'");
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createUser",
                    new String[] { "user_sp" }, txid))
                Util.log.info("Success: " + txid);

            // set user context
            Util.client.setUserContext(UserManagement.getOrCreateUser("user_sp"));

            Util.log.info("Create user `user_sp'");
            txid = new StringBuilder();
            if (Util.invoke("mapchannel", "mapcc", "createUser", new String[] { "user_sp" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Create AccountCL for the user `user_sp'");
            txid = new StringBuilder();
            if (Util.invoke("mapchannel", "mapcc", "createAccountCL",
                    new String[] { "user_sp", "Carol_Credit_0", "Credit", "zyp", "Nothing", "300", "zyp" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Split some of asset to sub-account");
            txid = new StringBuilder();
            if (Util.invoke("mapchannel", "mapcc", "splitAccountCL",
                    new String[] { "Carol_Credit_0", "user_sp", "Alice_Credit_1", "200" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Query parent account");
            res = Util.query("mapchannel", "mapcc", "queryAccount", new String[] { "Carol_Credit_0", "CL" });
            Util.log.info(res);

            Util.log.info("Query sub-account");
            res = Util.query("mapchannel", "mapcc", "queryAccount", new String[] { "Alice_Credit_1", "CL" });
            Util.log.info(res);
        } catch (Exception e) {
            Util.log.error(e.getMessage());
        }

        Util.log.info("End.");
        return;
    }
}
