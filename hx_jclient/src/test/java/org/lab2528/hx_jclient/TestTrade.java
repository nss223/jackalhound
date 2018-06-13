package org.lab2528.hx_jclient;

/**
 * Demo 2
 * <br>
 * trade different asset,
 * cross chain
 * 	
 * @author jack
 *
 */
public class TestTrade {
    public static void main(String[] args) {

        StringBuilder txid;
        String res;

        try {
            // set admin context
            Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));

            Util.log.info("Register user `traderA' and `traderB'");
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createUser", new String[] { "traderA" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createUser", new String[] { "traderB" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Register accounts of each user at each bank");
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "traderA", "A@a", "pointchannel", "points", "a_Bank" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "traderA", "A@b", "pointchannel", "points", "b_Bank" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "traderB", "B@a", "pointchannel", "points", "a_Bank" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "traderB", "B@b", "pointchannel", "points", "b_Bank" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Create the accounts");
            txid = new StringBuilder();
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "traderA", "A@a", "30", "a_Bank", "nothing" }, txid))
                Util.log.info("Success: " + txid);
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "traderA", "A@b", "40", "b_Bank", "nothing" }, txid))
                Util.log.info("Success: " + txid);
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "traderB", "B@a", "50", "a_Bank", "nothing" }, txid))
                Util.log.info("Success: " + txid);
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "traderB", "B@b", "60", "b_Bank", "nothing" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Extrade: @a A -> B = 7; @b B -> A = 6");
            txid = new StringBuilder();
            if (Util.invoke("pointchannel", "pointcc", "extrade",
                    new String[] { "A@a", "B@a", "7", "A@b", "B@b", "6" }, txid))
                Util.log.info("Success: " + txid);

            Util.log.info("Query account");
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "A@a", "all" });
            Util.log.info(res);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "B@a", "all" });
            Util.log.info(res);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "A@b", "all" });
            Util.log.info(res);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "B@b", "all" });
            Util.log.info(res);
        } catch (Exception e) {
            Util.log.error(e.getMessage());
        }

        try {
            Util.log.info("Register admin, the medium, and create the medium account");
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createUser", new String[] { "Admin" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "Admin", "Admin_a_Bank", "pointchannel", "points", "a_Bank" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "Admin", "Admin_a_Bank", "0", "a_Bank", "the medium of a_Bank" }, txid))
                Util.log.info("Success: " + txid);

            Util.client.setUserContext(UserManagement.getOrCreateUser("traderB"));
            Util.log.info("Request cross chain, especially at `a_Bank'");
            txid = new StringBuilder();
            if (Util.invoke("pointchannel", "pointcc", "crosstrade",
                    new String[] { "B@a", "Admin_a_Bank", "20", "another_point_channel" }, txid))
                Util.log.info("Success: " + txid);

            // after admin make sure the crosstrade is completed
            Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
            Util.log.info("Admin the medium register and create the corresponding account at `another_point_channel'");
            txid = new StringBuilder();
            if (Util.invoke("regchannel", "regcc", "createAccount",
                    new String[] { "traderB", "B@a_from_another", "pointchannel", "points", "a_Bank" }, txid))
                Util.log.info("Success: " + txid);
            txid = new StringBuilder();
            // the channel here should be another channel, i.e. `another_point_channel'
            // for simplicity we reuser `pointchannel'
            if (Util.invoke("pointchannel", "pointcc", "createAccount",
                    new String[] { "traderB", "B@a_from_another", "20", "point_a_Bank", "generated by cross request", "pointchannel", "a_Bank" }, txid))
                Util.log.info("Success: " + txid);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "B@a", "all" });
            Util.log.info(res);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "B@a_from_another", "all" });
            Util.log.info(res);
            res = Util.query("pointchannel", "pointcc", "queryAccount", new String[] { "Admin_a_Bank", "all" });
            Util.log.info(res);
        } catch (Exception e) {
            Util.log.error(e.getMessage());
        }
        Util.log.info("End.");
        return;
    }
}
