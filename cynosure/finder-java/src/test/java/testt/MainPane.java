package testt;

import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.value.ZkDataValue;
import com.iflytek.ccr.zkutil.ZkHelper;
import javafx.scene.Node;
import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.control.TextArea;
import javafx.scene.control.TextField;
import javafx.scene.layout.GridPane;
import javafx.scene.layout.VBox;
import utils.Md5Util;

import java.util.ArrayList;
import java.util.List;


/**
 * 主面板
 */
public class MainPane extends VBox {

    GridPane gridPane = new GridPane();
    TextField tfProject = new TextField("fff");
    TextField tfGroup = new TextField("GGG");
    TextField tfService = new TextField("GGG");
    TextField tfVersion = new TextField("GGG3");
    TextField tfZkAddress = new TextField("10.1.87.69:2183");
    TextArea tfRoute = new TextArea();
    TextArea tfConfig = new TextArea();


    Button btnCreateProject = new Button("Create");
    Button btnChangeRouteData = new Button("Change");
    Button btnChangeConfigData = new Button("Change");
    Button btnGetProviders = new Button("GetProviders");
    Button btnChangeZkAddr = new Button("Change");

    ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

    String pushId = "1234567890";

    String configData = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
    String  routeData = "[{\"id\":\"1\",\"consumer\":[\"c1\",\"c2\",\"c3\"],\"provider\":[\"p1\",\"p2\",\"p3\"],\"only\":\"Y\"}]";
    String instanceData = "1234567890{\"user\": {\"loadbalance\": \"loadbalance\",\"key1\": \"val\",\"key2\": \"val\"},\"sdk\": {\"is_valid\": true}}";

    List<Node> providerList = new ArrayList<>();

    public MainPane() {



        this.getChildren().addAll(gridPane);
        gridPane.addRow(0,new Label("ZkAddr"),tfZkAddress,btnChangeZkAddr);
        gridPane.addRow(1,new Label("Project"),tfProject);
        gridPane.addRow(2,new Label("Group"),tfGroup);
        gridPane.addRow(3,new Label("Service"),tfService);
        gridPane.addRow(4,new Label("Version"),tfVersion);
        gridPane.addRow(5,new Label("Route"),tfRoute,btnChangeRouteData);
        gridPane.addRow(6,new Label("ServerConfig"),tfConfig,btnChangeConfigData);
        getChildren().addAll(btnGetProviders);
        tfRoute.setPrefSize(300,100);
        tfConfig.setPrefSize(300,100);
        //设置route显示数据
        setPreData();
        setAutoChangeLine(tfRoute);
        addListener();

    }

    private void setPreData(){
        byte[] zkRouteData= zkHelper.getByteData(getBaseConfigPath() + "/route");
        ZkDataValue zkDataValue = ByteUtil.parseZkData(zkRouteData);
        if (zkDataValue == null || zkDataValue.getRealData() == null){
            tfRoute.setText(routeData);
        }
        else{
            tfRoute.setText(new String(zkDataValue.getRealData()));
        }
        tfConfig.setText(configData);

        byte[] zkConfigData= zkHelper.getByteData(getBaseConfigPath() + "/conf");
        ZkDataValue zkConfigValue = ByteUtil.parseZkData(zkConfigData);
        if (zkDataValue == null || zkConfigValue.getRealData() == null){
            tfConfig.setText(configData);
        }
        else{
            tfConfig.setText(new String(zkConfigValue.getRealData()));
        }
    }

    private void addListener(){
        btnChangeRouteData.setOnMouseClicked(event -> {

            boolean b = zkHelper.addOrUpdatePersistentNode(getBaseConfigPath()+"/route", ByteUtil.getZkBytes(tfRoute.getText().getBytes(), pushId));

        });

        btnChangeConfigData.setOnMouseClicked(event -> {
            boolean b = zkHelper.addOrUpdatePersistentNode(getBaseConfigPath()+"/conf", ByteUtil.getZkBytes(tfConfig.getText().getBytes(), pushId));

        });

        btnGetProviders.setOnMouseClicked(event -> {
            getProviders();
        });

        btnChangeZkAddr.setOnMouseClicked(event -> {
            zkHelper.closeClient();
            zkHelper = new ZkHelper(tfZkAddress.getText());
        });
    }

    private String  getBaseConfigPath(){
        String baseConfigPath;
        baseConfigPath = "/polaris/" +
                "service/" +
                Md5Util.getMD5(tfProject.getText()+tfGroup.getText())+"/" +
                tfService.getText() +"/"+
                tfVersion.getText();
        return baseConfigPath;
    }


    private void setAutoChangeLine(TextArea textArea){
        textArea.setStyle("ime-mode: active;\" class=\"active\"");
    }

    private void getProviders(){

        providerList.forEach(ch->getChildren().remove(ch));
        List<String> children = zkHelper.getChildren(getBaseConfigPath() + "/provider");
        for (String c:children){
            byte[] byteData = zkHelper.getByteData(getBaseConfigPath() + "/provider/" + c);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(byteData);
            String data;
            if (zkDataValue !=null && zkDataValue.getRealData() !=null) {
                data = new String(zkDataValue.getRealData());
            }
            else {
                data = "";
            }
            ProviderPane providerPane = new ProviderPane(getBaseConfigPath(),c,data,zkHelper);
            providerList.add(providerPane);
            getChildren().addAll(providerPane);
        }
    }
}
