package org.lab2528.hx_jclient;

import org.lab2528.hx_jclient.UserManagement;
import org.lab2528.hx_jclient.Util;

/**
 * Demo 4
 * <br>
 * data flow
 * 
 * @author jack
 *
 */
public class TestData {
	public static void main(String[] args) {
		
		StringBuilder txid;
		String res, user;

		try {			
			user = "Amy";
			Util.log.info("Set user context as `" + user + "'");
			Util.client.setUserContext(UserManagement.getOrCreateUser(user));
			Util.log.info("Commit data `id0'");
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "commit",
					new String[] { "id0", "db://id0", "key", "chash", "chash", "doc of id0" }, txid))
				Util.log.info("Success: " + txid);
			Util.log.info("Commit data `id1'");
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "commit",
					new String[] { "id1", "db://id1", "key", "chash", "chash", "doc of id1" }, txid))
				Util.log.info("Success: " + txid);
			
			user = "Bob";
			Util.log.info("Set user context as `" + user + "'");
			Util.client.setUserContext(UserManagement.getOrCreateUser(user));
			Util.log.info("Commit data `id3'");
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "commit",
					new String[] { "id3", "db://id3", "key", "chash", "chash", "doc of id3" }, txid))
				Util.log.info("Success: " + txid);
						
			Util.log.info("Checkout data");
			res = Util.query("datachannel", "datacc", "checkout", new String[] { "id3" });
			Util.log.info(res);
		} catch (Exception e) {
			Util.log.error(e.getMessage());
		}
		
		try {
			Util.log.info("Checkout others' data");
			res = Util.query("datachannel", "datacc", "checkout", new String[] { "id1" });
			Util.log.info(res);
		} catch (Exception e) {
			Util.log.error(e.getMessage());
		}
		
		try {
			Util.log.info("Share data without permission");
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "share",
					new String[] { "id1", "Bob" }, txid))
				Util.log.info("Success: " + txid);
		} catch (Exception e) {
			Util.log.error(e.getMessage());
		}
		
		try {
			Util.log.info("Share data by the owner");
			Util.client.setUserContext(UserManagement.getOrCreateUser("Amy"));
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "share",
					new String[] { "id1", "Bob" }, txid))
				Util.log.info("Success: " + txid);
			
			Util.client.setUserContext(UserManagement.getOrCreateUser("Bob"));
			Util.log.info("Checkout others' data after sharing");
			res = Util.query("datachannel", "datacc", "checkout", new String[] { "id1" });
			Util.log.info(res);
			
			Util.log.info("Make a new branch");
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "branch",
					new String[] { "id1", "id4", "db://id4", "key", "chash", "chash", "doc of id4" }, txid))
				Util.log.info("Success: " + txid);
			txid = new StringBuilder();
			if (Util.invoke("datachannel", "datacc", "branch",
					new String[] { "id4", "id5", "db://id5", "key", "chash", "chash", "doc of id5" }, txid))
				Util.log.info("Success: " + txid);
		} catch (Exception e) {
			Util.log.error(e.getMessage());
		}
		
		try {			
			Util.log.info("Check the trace of `id5'");
						res = Util.query("datachannel", "datacc", "trace", new String[] { "id5" });
			Util.log.info(res);
			
			Util.log.info("Check the history of `id1'");
			res = Util.query("datachannel", "datacc", "history", new String[] { "id1" });
			Util.log.info(res);

			Util.log.info("Collect data by owner `Bob'");
			res = Util.query("datachannel", "datacc", "queryByOwner", new String[] { });
			Util.log.info(res);

			Util.log.info("Collect data by create `Amy'");
			Util.client.setUserContext(UserManagement.getOrCreateUser("Amy"));
			res = Util.query("datachannel", "datacc", "queryByCreater", new String[] { });
			Util.log.info(res);
			
		} catch (Exception e) {
			Util.log.error(e.getMessage());
		}
		
		Util.log.info("End.");
		return;
	}
}
