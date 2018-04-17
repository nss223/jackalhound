package org.lab2528.multi;

import java.io.FileReader;
import java.util.Collection;
import java.util.Properties;
import java.util.concurrent.CompletableFuture;

import org.apache.log4j.Logger;
import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.EventHub;
import org.hyperledger.fabric.sdk.HFClient;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.QueryByChaincodeRequest;
import org.hyperledger.fabric.sdk.TransactionProposalRequest;
import org.hyperledger.fabric.sdk.BlockEvent.TransactionEvent;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.hyperledger.fabric.sdk.exception.TransactionException;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

/**
 * Basic functions
 * 
 * @author jack
 *
 */
public class Util {
	public static final Logger log = Logger.getLogger(Util.class);
	public static final Properties properties;
	public static HFClient client;
	public static Channel channel;
	
	private Util() {
	}

	static {
		properties = new Properties();
		try {
			
			// init config
			FileReader file = new FileReader("multi.properties");
			properties.load(file);
			file.close();
			
			// init HF client
			CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();
			client = HFClient.createNewInstance();
			client.setCryptoSuite(cryptoSuite);
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
	
	/**
	 * Get or initialize the channel
	 * the client must be bind with a user context first via `client.setUserContext(User)`.
	 * 
	 * @return Initialized channel
	 * @throws InvalidArgumentException
	 * @throws TransactionException
	 */
	public static Channel getChannel() throws InvalidArgumentException, TransactionException {
		if (null == channel)
			return getNewChannel();
		else
		return channel;
	}
	
	/**
	 * Force refresh the channel
	 * the client must be bind with a user context first via `client.setUserContext(User)`.
	 * 
	 * @return Initialized channel
	 * @throws InvalidArgumentException
	 * @throws TransactionException
	 */
	public static Channel getNewChannel() throws InvalidArgumentException, TransactionException
	{
        Peer peer = client.newPeer(
        		Util.properties.getProperty("peer"),
        		Util.properties.getProperty("peerEndpoint"));
        EventHub eventHub = client.newEventHub(
        		Util.properties.getProperty("eventHub"),
        		Util.properties.getProperty("eventHubEndpoint"));
        Orderer orderer = client.newOrderer(
        		Util.properties.getProperty("orderer"),
        		Util.properties.getProperty("ordererEndpoint"));
        channel = client.newChannel(
        		Util.properties.getProperty("channel"));
        
        channel.addPeer(peer);
        channel.addEventHub(eventHub);
        channel.addOrderer(orderer);
        channel.initialize();
        return channel;
	}
	
	/**
	 * Query chaincode, do not write
	 * 
	 * @param cc ChainCode
	 * @param fn Function
	 * @param args Args
	 * @return Response
	 * @throws InvalidArgumentException
	 * @throws ProposalException
	 * @throws TransactionException 
	 */
	public static String query(String cc, String fn, String[] args) throws InvalidArgumentException, ProposalException, TransactionException {
        QueryByChaincodeRequest qpr = client.newQueryProposalRequest();
        qpr.setChaincodeID(ChaincodeID.newBuilder().setName(cc).build());
        qpr.setFcn(fn);
        qpr.setArgs(args);
//        return getChannel().queryByChaincode(qpr);
        Collection<ProposalResponse> res = getChannel().queryByChaincode(qpr);
        return new String(res.iterator().next().getChaincodeActionResponsePayload());
//        for (ProposalResponse pres : res) {
//            String stringResponse = new String(pres.getChaincodeActionResponsePayload());
//            Util.log.info(stringResponse);
//        }
	}
	
	/**
	 * Invoke chaincode synclly
	 * 
	 * @param cc ChainCode
	 * @param fn Function
	 * @param args Args
	 * @param txid Return the transaction ID
	 * @return successful or not
	 * @throws ProposalException
	 * @throws InvalidArgumentException
	 * @throws TransactionException 
	 */
	public static boolean invoke(String cc, String fn, String[] args, StringBuilder txid) throws ProposalException, InvalidArgumentException, TransactionException {
    	TransactionProposalRequest tpr = client.newTransactionProposalRequest();
        tpr.setChaincodeID(ChaincodeID.newBuilder().setName(cc).build());
        tpr.setFcn(fn);
        tpr.setArgs(args);
        Collection<ProposalResponse> responses = getChannel().sendTransactionProposal(tpr);
        CompletableFuture<TransactionEvent> future = channel.sendTransaction(responses);
        TransactionEvent event = future.join();
        txid.setLength(0);
        txid.append(event.getTransactionID());
        return event.isValid();
	}
	
	/**
	 * Invoke chaincode asynclly
	 * 
	 * @param cc ChainCode
	 * @param fn Function
	 * @param args Args
	 * @return The future of transaction event
	 * @throws ProposalException
	 * @throws InvalidArgumentException
	 * @throws TransactionException 
	 */
	public CompletableFuture<TransactionEvent> invoke(String cc, String fn, String[] args) throws ProposalException, InvalidArgumentException, TransactionException {
		TransactionProposalRequest tpr = client.newTransactionProposalRequest();
        tpr.setChaincodeID(ChaincodeID.newBuilder().setName(cc).build());
        tpr.setFcn(fn);
        tpr.setArgs(args);
        Collection<ProposalResponse> responses = getChannel().sendTransactionProposal(tpr);
        return channel.sendTransaction(responses);
	}
}
