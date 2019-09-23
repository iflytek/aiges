package testt;

import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.control.TextArea;
import javafx.scene.layout.HBox;


public class ProviderPane extends HBox {

    Label lbaddress = new Label();
    TextArea instanceData = new TextArea();
    Button btnChange = new Button("change");

    public ProviderPane(String basePath,String address,String  data,ZkHelper zkHelper) {
        lbaddress.setText(address+"    ");
        instanceData.setText(data);
        instanceData.setPrefSize(300,150);
        getChildren().addAll(lbaddress,instanceData,btnChange);

        btnChange.setOnMouseClicked(event -> {
            zkHelper.update(basePath+"/provider/"+address,ByteUtil.getZkBytes(instanceData.getText().getBytes(),"1234567890"));

        });
    }

}
