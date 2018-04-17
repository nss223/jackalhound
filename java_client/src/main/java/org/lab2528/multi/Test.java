package org.lab2528.multi;

import java.security.KeyPair;
import java.util.Collection;

import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;

public class Test {
	public static void main(String[] args) {
		try {
			
			// invoke & query demo
			
			StringBuilder txid;
			String res;
			
			// set user context
			Util.client.setUserContext(UserManagement.getOrCreateUser("User1"));
			
			try {
				Util.log.info("Query user");
				res = Util.query("regcc", "queryUser", new String[] {"User1"});
				Util.log.info(res);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
		
			try
			{
				Util.log.info("Create account within user context");
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createAccount",
						new String[] { "User1", "AccountofBank", "mychannel", "points", "A_Bank" },
						txid))
					Util.log.info(txid);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}

			Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
			
			try {
				Util.log.info("Create account without creating user");
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createAccount",
						new String[] { "User1", "AccountofBank", "mychannel", "points", "A_Bank" },
						txid))
					Util.log.info(txid);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
			
			try {
				Util.log.info("Create user");
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createUser",
						new String[] { "User1"}, txid))
					Util.log.info("Success: " + txid);
				
				Util.log.info("Create account");
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createAccount",
						new String[] { "User1", "AccountofBank", "mychannel", "points", "A_Bank" },
						txid))
					Util.log.info("Success: " + txid);
				txid = new StringBuilder();
				if (Util.invoke("pointcc", "createAccount",
						new String[] { "User1", "AccountofBank", "300", "A_Bank", "Nothing" },
						txid))
					Util.log.info("Success: " + txid);

				Util.client.setUserContext(UserManagement.getOrCreateUser("User1"));
				Util.log.info("Query user");
				res = Util.query("regcc", "queryUser", new String[] {"User1"});
				Util.log.info(res);
				Util.log.info("Query account");
				res = Util.query("pointcc", "queryAccount", new String[] {"AccountofBank", "all"});
				Util.log.info(res);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}				
			
			try {
				Util.log.info("Set account by user");
				txid = new StringBuilder();
				if (Util.invoke("pointcc", "setAccount",
						new String[] { "User1", "AccountofBank", "400", "A_Bank", "Nothing" },
						txid))
					Util.log.info("Success: " + txid);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}				

			try {
				Util.log.info("Set account by admin");
				Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
				txid = new StringBuilder();
				if (Util.invoke("pointcc", "setAccount",
						new String[] { "User1", "AccountofBank", "400", "A_Bank", "Nothing" },
						txid))
					Util.log.info("Success: " + txid);

				Util.client.setUserContext(UserManagement.getOrCreateUser("User1"));
				Util.log.info("Query account after set");
				res = Util.query("pointcc", "queryAccount", new String[] {"AccountofBank", "all"});
				Util.log.info(res);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
			
			try {
				Util.log.info("Create asset mapping");
				Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createAccount",
						new String[] { "User1", "mychannel", "mychannel", "mapping", "ZYP" }, txid))
					Util.log.info("Success: " + txid);

				Util.log.info("Create asset user");
				txid = new StringBuilder();
				if (Util.invoke("mapcc", "createUser",
						new String[] { "User1" }, txid))
					Util.log.info("Success: " + txid);
				
				Util.log.info("Create asset");
				txid = new StringBuilder();
				if (Util.invoke("mapcc", "createAccount",
						new String[] {
								"User1",
								"a-car-BJ454852",
								"CAR",
								"Department of Motor Vehicles",
								"2018"
								}, txid))
					Util.log.info("Success: " + txid);
				
				Util.log.info("Create another asset");
				txid = new StringBuilder();
				if (Util.invoke("mapcc", "createAccount",
						new String[] {
								"User1",
								"a-house-YIHEYUAN load-5",
								"House",
								"Housing Authority",
								"Peking University"
								}, txid))
					Util.log.info("Success: " + txid);
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
			
			try {

				Util.client.setUserContext(UserManagement.getOrCreateUser("User2"));
				
				Util.log.info("Query asset account by another user");
				res = Util.query("mapcc", "queryAccount", new String[] { "a-car-BJ454852" });
				Util.log.info(res);

				Util.log.info("Query asset mapping by another user");
				res = Util.query("mapcc", "queryUser", new String[] {"User1" });
				Util.log.info(res);
				
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
			
			try {

				Util.client.setUserContext(UserManagement.getOrCreateUser("User1"));
				
				Util.log.info("Query asset account by himself");
				res = Util.query("mapcc", "queryAccount", new String[] { "a-car-BJ454852" });
				Util.log.info(res);

				Util.log.info("Query asset mapping by himself");
				res = Util.query("mapcc", "queryUser", new String[] {"User1" });
				Util.log.info(res);
				
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}
			
			try {

				Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
				
				Util.log.info("Query asset account by another admin");
				res = Util.query("mapcc", "queryAccount", new String[] { "a-car-BJ454852" });
				Util.log.info(res);

				Util.log.info("Query asset mapping by another admin");
				res = Util.query("mapcc", "queryUser", new String[] {"User1" });
				Util.log.info(res);
				
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}

			try {

				Util.client.setUserContext(UserManagement.getOrCreateUser("admin"));
				
				Util.log.info("Create another user and mapping");
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createUser",
						new String[] { "User2" }, txid))
					Util.log.info("Success: " + txid);
				txid = new StringBuilder();
				if (Util.invoke("regcc", "createAccount",
						new String[] { "User2", "mychannel", "mychannel", "mapping", "ZYP" }, txid))
					Util.log.info("Success: " + txid);
				txid = new StringBuilder();
				if (Util.invoke("mapcc", "createUser",
						new String[] { "User2" }, txid))
					Util.log.info("Success: " + txid);
				
				Util.log.info("User1 trade the car to User2");
				txid = new StringBuilder();
				if (Util.invoke("mapcc", "trade",
						new String[] { "User1", "User2", "a-car-BJ454852" }, txid))
					Util.log.info("Success: " + txid);
				
				Util.log.info("Check if User1 has the car now");
				res = Util.query("mapcc", "queryUser", new String[] { "User1" });
				Util.log.info(res);
				
				Util.log.info("Check if User2 has the car now");
				res = Util.query("mapcc", "queryUser", new String[] { "User2" });
				Util.log.info(res);

				Util.log.info("Check who has the car now");
				res = Util.query("mapcc", "queryAccount", new String[] { "a-car-BJ454852" });
				Util.log.info(res);
				
			} catch (Exception e) {
				Util.log.error(e.getMessage());
			}			
		
			Util.log.info("End.");
		} catch (InvalidArgumentException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		return;
	}
}
