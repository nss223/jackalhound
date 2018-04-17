package org.lab2528.multi;

import java.io.IOException;
import java.net.MalformedURLException;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.hyperledger.fabric_ca.sdk.exception.EnrollmentException;
import org.hyperledger.fabric_ca.sdk.exception.InvalidArgumentException;

/**
 * User load/enroll
 * 
 * @author jack
 *
 */
public class UserManagement {

	private static HFCAClient caClient = null;

	private static void initCAClient() throws MalformedURLException {
		// build CA client
		CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();
		caClient = HFCAClient.createNewInstance(
				Util.properties.getProperty("caEndpoint"), null);
		caClient.setCryptoSuite(cryptoSuite);
	}

	/**
	 * Load or create admin
	 * 
	 * @return An admin
	 * @throws ClassNotFoundException
	 * @throws IOException
	 * @throws EnrollmentException
	 * @throws InvalidArgumentException
	 */
	public static AppUser getOrCreateAdmin() throws ClassNotFoundException, IOException, EnrollmentException, InvalidArgumentException {
		AppUser admin = AppUser.load(
				Util.properties.getProperty("admin"));
		if (admin == null) {
			if (null == caClient)
				initCAClient();
			Enrollment adminEnrollment = caClient.enroll(
					Util.properties.getProperty("admin"),
					Util.properties.getProperty("admin_password"));
			admin = new AppUser(
					Util.properties.getProperty("admin"), 
					Util.properties.getProperty("affiliation"),
					Util.properties.getProperty("mspId"), adminEnrollment);
			admin.save();
		}
		return admin;
	}

	/**
	 * Get or create an user by name
	 * 
	 * @param userId User name
	 * @return An user
	 * @throws Exception
	 */
	public static AppUser getOrCreateUser(String userId) throws Exception {
		AppUser appUser = AppUser.load(userId);
		if (appUser == null) {
			if (null == caClient)
				initCAClient();
			RegistrationRequest rr = new RegistrationRequest(userId,
					Util.properties.getProperty("affiliation"));
			String enrollmentSecret = caClient.register(rr, getOrCreateAdmin());
			Enrollment enrollment = caClient.enroll(userId, enrollmentSecret);
			appUser = new AppUser(userId, 
					Util.properties.getProperty("affiliation"),
					Util.properties.getProperty("mspId"), enrollment);
			appUser.save();
		}
		return appUser;
	}
}
