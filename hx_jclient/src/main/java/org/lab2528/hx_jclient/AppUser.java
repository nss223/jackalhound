package org.lab2528.hx_jclient;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

import java.io.IOException;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.io.Serializable;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.PrivateKey;
import java.util.Set;

/**
 * Basic implementation of the {@link User} interface.
 *
 * @author jack
 *
 */
public class AppUser implements User, Serializable {

    private static final long serialVersionUID = 1L;
    private String name;
    private Set<String> roles;
    private String account;
    private String affiliation;
    private Enrollment enrollment;
    private String mspId;

    public AppUser(String name, String affiliation, String mspId, Enrollment enrollment) {
        this.name = name;
        this.affiliation = affiliation;
        this.enrollment = enrollment;
        this.mspId = mspId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Set<String> getRoles() {
        return roles;
    }

    public void setRoles(Set<String> roles) {
        this.roles = roles;
    }

    public String getAccount() {
        return account;
    }

    public void setAccount(String account) {
        this.account = account;
    }

    public String getAffiliation() {
        return affiliation;
    }

    public void setAffiliation(String affiliation) {
        this.affiliation = affiliation;
    }

    public Enrollment getEnrollment() {
        return enrollment;
    }

    public void setEnrollment(Enrollment enrollment) {
        this.enrollment = enrollment;
    }

    public String getMspId() {
        return mspId;
    }

    public void setMspId(String mspId) {
        this.mspId = mspId;
    }

    @Override
    public String toString() {
        return "AppUser: " + name + "\n" + enrollment.getCert();
    }


    /**
     * @return Cert in X509
     */
    public String getCert() {
        return this.enrollment.getCert();
    }

    /**
     * @return Key in java.security.PrivateKey
     */
    public PrivateKey getKey() {
        return this.enrollment.getKey();
    }

    /**
     * Load user from local binary; if failed, then create
     *
     * @param name User name
     * @return The user
     * @throws IOException
     * @throws ClassNotFoundException
     */
    public static AppUser load(String name) throws IOException, ClassNotFoundException {
        if (Files.exists(Paths.get(name + ".jso"))) {
             ObjectInputStream decoder = new ObjectInputStream(
                     Files.newInputStream(Paths.get(name + ".jso")));
             return (AppUser)decoder.readObject();
        } else {
            return null;
        }
    }

    /**
     * Save user obj to binary
     *
     * @throws IOException
     */
    public void save() throws IOException {
        ObjectOutputStream oos = new ObjectOutputStream(
                Files.newOutputStream(Paths.get(this.name + ".jso")));
        oos.writeObject(this);
    }
}