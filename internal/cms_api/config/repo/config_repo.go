package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"im/internal/cms_api/config/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"sort"
	"time"

	"gorm.io/gorm/clause"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const ConfigMenuConfTimeStamp = "CONFIG_MENU_CONF_TIMESTAMP"
const ConfigDiscoverIsOpen = "CONFIG_DISCOVER_IS_OPEN"
const ConfigGoogleCodeIsOpen = "CONFIG_GOOGLE_CODE_IS_OPEN"
const ConfigCmsSiteName = "CONFIG_CMS_SITE_NAME"
const ConfigCmsLoginIcon = "CONFIG_CMS_LOGIN_ICON"
const ConfigCmsLoginBackend = "CONFIG_CMS_LOGIN_BACKEND"
const ConfigCmsPageIcon = "CONFIG_CMS_PAGE_ICON"
const ConfigCmsMenuIcon = "CONFIG_CMS_MENU_ICON"
const ConfigJPushAppKey = "CONFIG_JUSH_APP_KEY"
const ConfigJPushMasterSecret = "CONFIG_JUSH_MASTER_SECRET"
const ConfigFeihuAppKey = "CONFIG_FEIHU_APP_KEY"
const ConfigFeihuAppSecret = "CONFIG_FEIHU_APP_SECRET"
const ConfigDepositHTML = "CONFIG_DEPOSIT_HTML"
const ConfigDepositURL = "CONFIG_DEPOSIT_URL"
const ConfigDepositSWITCH = "CONFIG_DEPOSIT_SWITCH"
const ConfigUserAgreement = "CONFIG_USER_AGREEMENT"
const ConfigPrivacyPolicy = "CONFIG_PROVACY_POLICY"
const ConfigAboutUs = "CONFIG_ABOUT_US"
const ConfigIPWhiteIsOpen = "CONFIG_IP_WHITE_IS_OPEN"
const ConfigPrivilegeUserFreeze = "CONFIG_PRIVILEGE_USER_FREEZE"
const ConfigRegisterConfig = "register_config"
const ConfigLoginConfig = "login_config"
const ConfigAnnouncementConfig = "announcement_config"
const ConfigSignConfig = "sign_log_config"
const ConfigParameterConfig = "system_config"
const ConfigWithdrawConfig = "withdraw_config"
const ConfigDefaultIsOpen = "CONFIG_DEFAULT_IS_OPEN"
const ConfigDefaultGroupIndex = "CONFIG_DEFAULT_GROUP_INDEX"
const ConfigDefaultFriendIndex = "CONFIG_DEFAULT_FRIEND_INDEX"

const DefaultAboutUs string = "<p><strong>关于我们</strong></p><p>&nbsp;</p><p>IM通讯成立于2022年3月，是IM领域领先的企业级软件服务提供商。旗下产品线主要包括IM通讯即时通讯云、IM通讯实时音视频、IM通讯客服系统，以及企业级IM通讯机器人，是行业内较早覆盖通讯、客服、智能机器人的一体化通讯产品技术服务公司。</p><p>&nbsp;</p><p>IM即时通讯支持以下服务内容：</p><p>--支持即时通讯，单聊、群聊、聊天室</p><p>--支持实时音视频，单人&amp;多人音视频通话</p><p>--支持多平台，一次开发覆盖多平台</p><p>--支持文字、语音、视频、图片、表情、自定义等消息</p><p>--其他定制化开发内容</p>"
const DefaultPrivacyPolicy string = "<p><strong>隐私政策</strong></p><p>IM通讯非常重视保护您的隐私。</p><p>为方便您登录、使用相关服务，以及为您提供更个性化的用户体验和服务，您在使用我们的服务时，我们可能会收集和使用您的相关信息。我们希望通过本隐私声明向您说明，在使用IM通讯服务（统称“本服务”）时，我们如何收集、使用、储存和披露您的信息，以及我们为您提供的访问、更新、控制和保护这些信息的方式。本隐私声明与您所使用的IM通讯服务息息相关，希望您仔细阅读。&nbsp;您使用我们的服务，即意味着您已经同意我们按照本隐私声明收集、使用、储存和披露您的相关信息，以及向您提供的控制和保护措施。</p><p>&nbsp;</p><p><strong>一、我们收集哪些您的个人信息</strong></p><p>就本政策而言，“个人信息”是指以电子或者其他方式记录的能够单独或者与其他信息结合识别特定自然人身份或者反映自然人活动情况的各种信息。</p><p>为了向您提供服务、保障服务的正常运行、改进和优化我们的服务以及保障帐号安全，我们会主动向您收集为提供服务所必需的个人信息。IM通讯会按照如下方式收集你在注册、使用服务时主动提供、授权提供或因为使用服务而产生的信息：</p><p><strong>1.1当您注册账户时：</strong></p><p>在您注册使用本服务时，您将向我们提供您的账号，并设置您的密码。您也可以选择添加其他信息来完善您的账户，例如上传您的照片。为完成支付，还需要设置支付密码。</p><p><strong>1.2当您使用我们的产品时：</strong></p><p>在您使用本服务期间，我们将收集为提供持续服务和保证服务质量所需的如下信息：</p><p><strong>a.与您设备相关的信息：</strong>您的设备随机ID、设备硬件类型、设备型号、设备系统类型；</p><p><strong>b.与您网络状况相关的信息：</strong>您的设备网络类型；</p><p><strong>c.您与我们产品交互相关的信息：</strong>收发端用户ID、APPID（应用包名）、登录用时、SDK版本；</p><p><strong>d.您的聊天信息：</strong>我们为企业客户提供IM服务及相关配套服务，若您作为企业客户的终端用户使用我们的该项功能时或者以个人用户名义使用我们的Demo时，我们将为实现通讯业务功能收集、传输、储存您在使用服务过程中产生的文字通讯信息及音视频流信息，这些信息可能包括您或者/及与您进行互动的个人向我们提供的您的个人信息（包括您发送的文字、图片、音频、视频、您的图像、声音、肖像以及您交流的内容）；</p><p><strong>e.当您使用位置消息时：</strong>将收集您的GPS定位信息、经纬度；</p><p><strong>1.3当您使用我们的APP时：</strong></p><p>我们将获取您的运行中的进程，确保 SDK&nbsp;只在主进程完成初始化。</p><p><strong>1.4当您提出服务咨询时：</strong></p><p>当您在使用我们产品过程中向我们反馈服务质量、申请技术支持、或反馈产品优化建议时，或者通过我们的官网、邮箱、电话向我们提出申请、要求、诉求、投诉时，我们将收集您的姓名、邮箱、通讯地址、电话，以及您申请、要求、诉求及投诉的内容。为加强沟通，若您主动提供，我们也将收集您的其他相关信息。</p><p>&nbsp;</p><p><strong>二、我们如何存储这些信息</strong></p><p><strong>2.1&nbsp;信息存储的地点</strong></p><p>我们会按照法律法规规定，进行手机和存储必要的信息。</p><p><strong>2.2&nbsp;信息存储的期限</strong></p><p>通常，我们仅在为您提供服务期间保留您的信息，保留时间不会超过满足相关使用目的所必须的时间。&nbsp;</p><p>但在下列情况下，且仅出于下列情况相关的目的，我们有可能需要较长时间保留您的信息或部分信息：&nbsp;</p><p>•&nbsp;遵守适用的法律法规等有关规定。&nbsp;</p><p>•&nbsp;遵守法院判决、裁定或其他法律程序的要求。</p><p>•&nbsp;遵守相关政府机关或其他有权机关的要求。</p><p>•&nbsp;我们有理由确信需遵守法律法规等有关规定。</p><p>•&nbsp;为执行相关服务协议或本隐私声明、维护社会公共利益、处理投诉/纠纷，保护我们的客户、我们或我们的关联公司、其他用户或雇员的人身和财产安全或合法权益所合理必需的用途。</p><p>&nbsp;</p><p><strong>三、我们如何保护这些信息</strong></p><p>3.1&nbsp;我们努力为用户的信息安全提供保障，以防止信息的丢失、不当使用、未经授权访问或披露。</p><p>3.2&nbsp;我们将在合理的安全水平内使用各种安全保护措施以保障信息的安全。例如，我们将通过服务器多备份、密码加密等安全措施，防止信息泄&nbsp;露、毁损、丢失。</p><p>3.3&nbsp;我们建立严格的管理制度和流程以保障信息的安全。例如，我们严格限制访问信息的人员范围，并进行审计，要求&nbsp;他们遵守保密义务。</p><p>3.4&nbsp;若发生个人信息泄露等安全事件，我们会启动应急预案，阻止安全事件扩大，按照《国家网络安全事件应急预案》&nbsp;等有关规定及时上报，并以发送邮件、推送通知、公告等形式告知您相关情况，并向您给出安全建议。</p><p>3.5&nbsp;我们重视信息安全合规工作，并通过众多国际和国内的安全认证，以业界先进的解决方案充分保障您的信息安全。&nbsp;</p><p>我们会尽力保护你的个人信息。我们也请你理解，任何安全措施都无法做到无懈可击。</p><p>&nbsp;</p><p><strong>四、我们如何使用这些信息</strong></p><p>为了向您提供更加优质、便捷、安全的服务，在符合相关法律法规的前提下，我们可能将收集的信息用作以下用途：&nbsp;</p><p>•&nbsp;向您提供服务。</p><p>•&nbsp;满足您的个性化需求。例如，语言设定、位置设定、个性化的帮助服务和指示，或对您和其他用户作出其他方面的&nbsp;回应。&nbsp;</p><p>•&nbsp;服务优化和开发。例如，我们会根据IM通讯系统响应您的需求时产生的信息，优化我们的服务。&nbsp;</p><p>•&nbsp;保护IM通讯、IM通讯用户和IM通讯的合作伙伴。例如，我们会将您的信息用于身份验证、安全防范、投诉处理、纠纷协调、诈骗监测等用途。例&nbsp;</p><p>•&nbsp;向您提供与您更加相关的服务。例如，向您提供您可能感兴趣的类似功能或服务等。</p><p>•&nbsp;邀请您参与有关我们产品和服务的调查。</p><p>•&nbsp;其他可能需要使用收集的信息的相关场景，如果使用场景与初始场景无合理的关联性，我们会在使用信息前重新征得您的同意。</p><p>&nbsp;</p><p><strong>五、信息的分享和对外提供</strong></p><p>我们不会与任何无关第三方分享您的信息，但以下情况除外：</p><p>5.1&nbsp;获取你的明确同意：经你事先同意，我们可能与第三方分享你的个人信息；</p><p>5.2&nbsp;为实现外部处理的目的，我们可能会与关联公司或其他第三方合作伙伴（第三方服务供应商、承包商、代理等）分享你的个人信息，让他们按照我们的说明、隐私政策以及其他相关的保密和安全措施来为我们处理上述信息，并用于向你提供我们的服务，实现“我们如何使用信息”部分所述目的。如我们与任何上述第三方分享您的信息，我们将努力确保第三方在使用您的信息时遵守本声明及我们要求其遵守的其他适当的保密和安全措施。&nbsp;</p><p>5.3&nbsp;随着我们业务的持续发展，我们以及我们的关联公司有可能进行合并、收购、资产转让或类似的交易，您的信息有可能作为此类交易的一部分而被转移。我们将遵守相关法律法规的要求，在转移前通知您，确保信息在转移时的&nbsp;机密性，以及变更后继续履行相应责任和义务。</p><p>5.4&nbsp;我们还可能因以下原因而披露您的信息：</p><p>•遵守适用的法律法规等有关规定。&nbsp;</p><p>•遵守法院判决、裁定或其他法律程序的规定。</p><p>•遵守相关政府机关或其他有权机关的要求。&nbsp;</p><p>•我们有理由确信需遵守法律法规等有关规定。&nbsp;</p><p>•为执行相关服务协议或本隐私声明、维护社会公共利益、处理投诉/纠纷，保护我们的客户、我们或我们的关联公司、其他用户或雇员的人身和财产安全或合法权益所合理必需的用途。&nbsp;</p><p>•经过您合法授权的情形。&nbsp;</p><p>如我们因上述原因而披露您的信息，我们将在遵守法律法规相关规定及本声明的基础上及时告知您。</p><p>&nbsp;</p><p><strong>六、你如何访问及管理个人信息</strong></p><p>6.1&nbsp;您可以在使用我们服务的过程中，访问、修改和删除您的相关信息。您访问、修改和删除信息的方式将取决于您使用的具体服务。&nbsp;</p><p>6.2&nbsp;如您发现我们违反法律法规的规定或者双方的约定收集、使用您的信息，您可以要求我们删除。如您发现我们收集、存储的您的信息有错误的，且无法自行更正的，您也可以要求我们更正。&nbsp;</p><p>6.3&nbsp;在访问、修改和删除相关信息时，我们可能会要求您进行身份验证，以保障账户安全。请您理解，由于技术所限、基于法律法规要求，您的某些请求可能无法进行响应。&nbsp;</p><p>6.4&nbsp;基于法律法规要求、保障信息安全等正当事由，您的部分信息可能无法访问、修改和删除。&nbsp;</p><p>6.5&nbsp;您可以自行选择撤回对某些非基本功能或服务对您信息处理的同意，并通过联系客服的方式选择注销IM通讯账号。&nbsp;</p><p>6.6&nbsp;如您对上述权利实现存在疑问，请联系官方客服。</p><p>&nbsp;</p><p><strong>七、本隐私指引的变更</strong></p><p>我们可能适时修订本隐私声明的条款，该等修订构成本《隐私声明》的一部分。如可能造成您在本隐私声明下权利的实质减少或扩大收集、使用信息的范围等重要规则变更时，我们将在修订生效前通过在主页上显著位置提示或向您发送电子邮件或以其他方式通知您。在该种情况下，若您继续使用我们的服务，即表示同意受经修订的本隐私声明的约束。</p><p>&nbsp;</p><p><strong>八、与我们联系</strong></p><p>如果您对本服务条款有任何问题或建议或有其他事情需要联系我们，请联系官方客服。</p>"
const DefaultUserAgreement string = "<p><strong>用户协议</strong></p><p>&nbsp;</p><p><strong>一、接受服务条款</strong></p><p>欢迎您使用我们提供的IM通讯服务（以下简称：IM服务），请仔细阅读下列内容并在明确理解的前提下注册：IM服务的所有权和运营权归IM通讯所有。为获得IM通讯所提供的相关服务，您需要同意本条款内容并按照页面上的提示完成注册程序。如果您在注册过程中勾选“同意遵守【IM服务使用协议】”选项即表示完全接受本条款的所有内容，并且在知情的基础上明确同意按条款的规定处理和使用。如果您不接受本条款，请不要安装、使用、注册或以其他方式访问我们的服务。<br>“我们”或“本公司”：即软件所有人IM通讯，为IM通讯的开发商，依法独立拥有IM通讯产品著作权。</p><p>“本软件”：IM通讯。</p><p>“您”：使用本软件的任何个人、公司或其他组织，无论是否盈利、何种用途（包括以学习和研究为目的）。</p><p>“本协议”、 “本条款”：即【用户协议】，本协议适用且仅适用于IM通讯，IM通讯拥有在法律允许范围内对本协议的最终解释权。</p><p>&nbsp;</p><p><strong>二、资格和授权</strong></p><p>1、为了使用IM服务，您需要先同意本协议及《用户隐私协议》，成为我们的网站注册账户，才能使用我们的服务。为了注册账户的目的，您需要输入账号、设置密码并填写相关必要信息以完成注册。请您务必提供真实、准确的注册信息，以便能及时接收我们发送的通知消息，我们可以利用这些信息来通知您最新的产品更新和市场活动。若注册信息有任何变动，请您及时更新；如因您提供的信息存不正确、不真实、不完整或未及时更新等导致给您带来损失、费用增加及其他不利后果将由您自行承担。</p><p>2、您承诺您是具备完全民事权利能力和完全民事行为能力的自然人、法人或其他组织，您应对使用该账户及密码进行的一切操作及言论承担全部责任。</p><p>3、您应确保您在每个上网时段结束时，以正确步骤离开网站。本公司不能也不会对因您未能遵守本款规定而发生的任何损失或损毁负责。您应妥善保管您的账户和密码，并对账户、密码等信息采取必要、有效的保密和安全保护措施，如设置高强度密码、定期更换等；只有您本人可以使用您的账户，除非有法律规定或司法裁定，或经本公司同意外，该账户不可转让、不可赠与、不可继承、不可通过出借、出租等任何方式提供给其他第三方使用（与帐户相关的财产权益除外）。因他人未经授权使用您的用户名和密码造成的损失由您自行承担或依法由擅用人承担，本网站尽力保证您的帐户安全但不保证您帐户安全的绝对性。</p><p>4、若您是企业客户并将账户授权您的员工管理的，您应自行建立健全内部管理制度，做好权限管控，并且在遇到人员变动时及时完成交接和账户安全保护。对于因您未妥善保管账户或密码（包括但不限于向第三者透露账户和密码及其他注册资料、多人共享同一个账户、安装非法或来源不明的程序、将已经登陆IM服务的账户的设备提供给他人使用等）导致泄露或遭他人非法使用，将由您自行承担相应法律后果。</p><p>5、如您发现有他人冒用或盗用您的账户及密码或任何其他未经合法授权之情形时，应立即联系我们，并提供必要资料（如客户资料、情况说明、证明材料及诉求等）。我们收到您的有效通知并核实身份后，会依据法律法规及服务规则进行处理。若因您的原因（如提供的资料存在瑕疵），导致我们无法核实您的身份或无法判断您的诉求等进而未能及时处理而给您带来的损失，应由您自行承担。同时，您应理解我们对您的请求采取行动需要合理期限，对于我们采取措施之前已执行的指令、您已经产生的损失以及采取措施后因不可归责于我们的原因导致的损失，我们不承担任何责任。</p><p>&nbsp;</p><p><strong>三、服务内容</strong></p><p>1、本条款中的“产品/服务”指我们向您提供IM通讯官方的产品及其他经双方协商选定的产品/服务，以及相关的技术及网络支持服务。具体产品/服务内容以您所订购的产品/服务内容为准。</p><p>2、IM服务的具体内容由IM通讯根据实际情况提供给您，我们保留随时修改本条款的权利，并会在修改后于本页面公布修改后的文本，因此，请经常查看本页。若您不同意修改后的文本，请立即停止使用；如果您继续使用我们的服务，就视同接受我们对本条款的修改。</p><p>&nbsp;</p><p>3、若您的应用存在滥用系统资源，包括但不限于：因您程序逻辑、设计上的缺陷或者您系统被攻击导致异常频率和数量的接口调用，引发大量不合理的系统资源消耗行为；或者脱离正常业务模型大量发送消息或者持续连接服务器等，本公司保留对您强制停止提供服务及冻结账户的权利。</p><p>&nbsp;</p><p><strong>四、使用限制</strong></p><p>1、对于您利用IM服务所发布的信息，IM通讯保留依据相关法律法规对其通讯的信息进行关键词过滤的权利，如发现您发送内容明确存在违反相关法律法规的，IM通讯有权作出包括但不限于劝阻、拦截、直至向有关部门举报等行为。但这并不表示IM通讯对您所发送的内容有过滤或审核的义务，也没有任何审查、审核、监督的责任或其他连带责任。您需自行对发送该等信息的行为承担一切责任，与IM通讯无关，必要时您还需承担由此给IM通讯造成的相关损失。</p><p>2、您不得对IM服务（包括但不限于IM服务提供的任何SDK、本网站）进行出售、转移、转让、出租、租赁、再许可、修改、二次创作、反编译或反向工程，或其他提取源代码的尝试以及实施任何涉嫌侵害我们合法权益的行为，除非我们给予您明确的书面认可。</p><p>3、您在使用我们提供的IM服务时，必须遵循以下原则：<br>　　　　 &nbsp; 不得危害国家安全、泄露国家秘密，不得侵犯国家社会集体的和公民的合法权益，不得制作、复制、查阅和传播下列信息：<br>　　　　　　违反宪法确定的基本原则的；<br>　　　　　　危害国家安全，泄漏国家机密，颠覆国家政权，破坏国家统一的；<br>　　　　　　损害国家荣誉和利益的；<br>　　　　　　煽动民族仇恨、民族歧视，破坏民族团结的；<br>　　　　　　破坏国家宗教政策，宣扬邪教和封建迷信的；<br>　　　　　　散布谣言，扰乱社会秩序，破坏社会稳定的；<br>　　　　　　散布淫秽、色情、赌博、暴力、恐怖或者教唆犯罪的；<br>　　　　　　侮辱或者诽谤他人，侵害他人合法权益的；<br>　　　　　　煽动非法集会、结社、游行、示威、聚众扰乱社会秩序的；<br>　　　　　　以非法民间组织名义活动的；<br>　　　　　　含有法律、行政法规禁止的其它内容的。<br>　　　　　　不得用任何不正当手段损害我们及其他用户的利益及声誉。<br>　　　　　　违反上述规定的，我们有权终止对您进行服务，并协助有关行政机关或司法机关等进行追索和查处。</p><p>4、非经我们开发或我们授权开发并正式发布的其它任何由本IM服务衍生的软件均属非法，下载、安装、使用此类衍生软件，将可能导致不可预知的风险，建议您不要轻易下载、安装、使用，由此产生的一切法律责任与纠纷一概与我们无关。</p><p>5、如将本软件应用于商业用途，您必须获得我们的商业授权。获得商业授权之后，您可以将本软件应用于授权指定的商业用途，同时依据所购买的授权类型中确定的技术支持期限、技术支持方式和技术支持内容，自购买时刻起，在技术支持期限内拥有通过指定的方式获得指定范围内的技术支持服务。商业授权用户享有反映和提出意见的权利，相关意见将被作为重要考虑内容，但没有一定被采纳的承诺或保证。未获商业授权之前，任何单位或个人不得将本软件用于商业用途（包括但不限于企业网站或软件、政府单位网站或软件、经营性网站、以盈利为目的的网站或软件），否则您应承担相应法律责任，包括但不限于赔偿我们因维护权益而产生的律师费、差旅费、鉴定费、公证费、仲裁费等所有追偿费用。购买商业授权请与商务联系。</p><p>&nbsp;</p><p><strong>五、保密条款</strong></p><p>1、保密资料指由一方向另一方披露的所有技术及非技术信息(包括但不限于产品资料，产品计划，价格，财务及营销规划，业务战略，客户信息，客户数据，研发，软件硬件，API应用数据接口，技术说明，设计，特殊公式，特殊算法等)。</p><p>2、本协议任何一方同意对获悉的对方之上述保密资料予以保密，并严格限制接触上述保密信息的员工遵守本条之保密义务。除非国家机关依法强制要求或上述保密资料已经进入公有领域外，接受保密资料的一方不得对外披露。</p><p>3、本协议双方明确认可各自用户信息和业务数据等是各自的重要资产及重点保密信息。本协议双方同意尽最大的努力（至少不低于对待自己的保密信息的谨慎）保护上述保密信息等不被披露。一旦发现有上述保密信息泄露事件，双方应合作采取一切合理措施避免或者减轻损害后果的产生。</p><p>4、尽管有前述约定，符合下列情形之一的，相关信息不被视为保密信息：<br>　　　　4.1&nbsp;接收方在披露方向其披露之前已经通过合法的渠道或方式持有的信息。<br>　　　　4.2&nbsp;该信息已经属于公知领域，或该信息在非因接收方违反本协议的情况下而被公开。<br>　　　　4.3&nbsp;接收方合法自其他有权披露资料的第三方处知悉且不负有保密义务的信息。<br>　　　　4.4&nbsp;由接收方不使用或不参考任何披露方的保密信息而独立获得或开发的。</p><p>5、如果接收方基于法律法规或监管机关的要求，需要依法披露披露方的保密信息的，不视为违反本协议，但接收方应当在法律许可的范围内尽快通知披露方，同时，接收方应当努力帮助披露方有效限制该保密信息的披露范围，保护披露方合法权益。</p><p>6、双方保密义务在本协议有效期限内及期限届满后持续有效，直至相关信息不再具有保密意义。</p><p>7、一旦发生保密信息泄露事件，双方应合作采取一切合理措施避免或者减轻损害后果的产生；如因接收方违反保密义务给披露方造成损失的，接收方应赔偿因此给披露方造成的直接经济损失。</p><p>8、您同意在不披露您保密信息的前提下，我们可以就您使用IM服务的情况作为使用范例或成功案例用于IM服务自身及/或业务的宣传与推广，在这类宣传推广中，IM服务有权免费在全球范围内使用您的名称、企业标识、相关logo&nbsp;等；除前述情况外，我们不会未经您书面许可使用前述信息。</p><p>&nbsp;</p><p><strong>六、隐私保护</strong></p><p>我们将按照法律法规要求，采取安全保护措施，保护您的用户信息安全可控。具体详见《用户隐私协议》。</p><p>&nbsp;</p><p><strong>七、知识产权</strong></p><p>1、IM服务由IM通讯研发。除非另有说明，IM服务的一切版权等知识产权，以及与IM服务相关的所有信息内容，包括但不限于著作、图片、档案、资讯、资料、架构、页面设计、本公司网站Logo、“IM通讯”等文字、图形及其组合，以及网站的其他标识、徽记、服务的名称、技术文档等，均由本公司依法拥有其知识产权。</p><p>2、除IM通讯或第三方明示同意外，您无权复制、传播、转让、许可或提供他人使用上述资源，否则必须承担相应的责任。</p><p>&nbsp;</p><p><strong>八、免责声明</strong></p><p>1、如发生下述情形，本公司不承担任何法律责任：<br>　　　　1.1依据法律规定或相关政府部门的要求提供您的个人信息；<br>　　　　1.2由于您的使用不当而导致任何个人信息的泄露；<br>　　　　1.3任何由于黑客攻击，电脑病毒的侵入，非法内容信息、骚扰信息的屏蔽，政府管制以及其他任何网络、技术、通信线路、信息安全管理措施等原因造成的服务中断、受阻等不能满足用户要求的情形；<br>　　　　1.4用户因第三方如运营商的通讯线路故障、技术问题、网络、电脑故障、系统不稳定及其他因不可抗力造成的损失的情形；<br>　　　　1.5使用IM服务可能存在的来自他人匿名或冒名的含有威胁、诽谤、令人反感或非法内容的信息而招致的风险。</p><p>2、本公司的服务按“现状”和“现有”提供，不提供任何形式的明示或暗示的担保，并且本公司明确免除任何担保和条件，包括但不限于任何暗示担保，适用于特定用途，合法性、可用性、安全性、所有权和/或非侵权性。您对网站和本公司服务的使用由您自行决定并承担风险，并且您将对因使用而造成的任何损害承担全部责任，包括但不限于对您的计算机系统造成的任何损害，或数据的丢失或损坏。</p><p>3、本公司不对与本服务条款有关或由本服务条款引起的任何间接的、惩罚性的、特殊的、派生的损失承担赔偿责任。</p><p>&nbsp;</p><p><strong>九、赔偿和豁免</strong></p><p>1、您同意对本公司及其子公司、高级职员、董事、代理商、服务提供商、合作伙伴和员工因下列因素产生的任何第三方索赔、要求或指控，以及所有相关损失、损害、责任、成本和开支（包括律师费）进行赔偿、抗辩及保护其不受损害：<br>　　　（1）您对本网站或IM服务的使用；<br>　　　（2）您的产品的用户对IM服务的使用；<br>　　　（3）您的产品，包括您产品上的任何内容、服务或广告，或您与IM服务合并使用的内容，服务或广告；<br>　　　（4）任何因未经授权使用IM服务而引起或与之相关的侵犯版权、诽谤、侵犯隐私权或公开权的任何索赔；<br>　　　（5）您违反本条款中包含的任何陈述、保证或契约。</p><p>2、上述赔偿条款应是对这些条款中规定的任何其他赔偿义务的增补而非替代。</p><p>&nbsp;</p><p><strong>十、出口合规</strong></p><p>双方承诺遵守本协议所适用的联合国、中国、美国以及其他国家出口管制法律法规。您承诺不会将本公司提供的产品或服务用于适用的出口管制法律法规禁止的用途。非经相关主管机关许可，您及您授权使用本公司产品或服务的其他个人或实体不会通过本公司产品或服务向适用的出口管制法律法规禁止的实体或个人提供受管控的技术、软件或服务。</p><p>此外，您同意在法律允许的最大范围内赔偿本公司您违反此规定可能产生的任何罚款或罚金。此出口管制条款在终止或取消本条款后仍然有效。</p><p>&nbsp;</p><p><strong>十一、其他条款</strong></p><p>1、本服务协议所定的任何条款的部分或全部无效的，不影响其它条款的效力。</p><p>2、本服务协议的解释、效力及纠纷的解决。若您和我们之间发生任何纠纷或争议，首先应友好协商解决，协商不成的，您同意将纠纷或争议提交本公司所在地有管辖权的法院诉讼解决。</p><p>3、本条款项下之保密条款、知识产权条款、法律适用及争议解决条款等内容，不因本条款的终止而失效。</p><p>4、若双方之间另有签署协议，且与网络页面点击确认的本条款存在不一致之处，以双方签署协议内容为准。</p><p>5、本协议自上线之日起生效。</p>"

var ConfigRepo = new(configRepo)

type configRepo struct{}

func (r *configRepo) MenuGetConfigTime() (timeString string, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigMenuConfTimeStamp).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		t := time.Now().UnixNano()
		logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("time", fmt.Sprintf("%d", t)))
		timeString = util.Int64ToString(t)
		logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("timeString", timeString))
		config.Name = ConfigMenuConfTimeStamp
		config.Value = timeString
		db.DB.Save(&config)
	} else {
		timeString = config.Value
	}
	return
}

func (r *configRepo) MenuUpdateConfigTime(timeString string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigMenuConfTimeStamp).Update("value", timeString).Error
	return
}

func (r *configRepo) GetDiscoverIsOpen() (status int, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDiscoverIsOpen).First(&config).Error
	status = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetDiscoverIsOpen(status int) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDiscoverIsOpen).Update("value", status).Error
	return
}

func (r *configRepo) GetGoogleCodeIsOpen() (status int, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigGoogleCodeIsOpen).First(&config).Error
	status = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetGoogleCodeIsOpen(status int) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigGoogleCodeIsOpen).Update("value", status).Error
	return
}

func (r *configRepo) GetIPWhiteIsOpen() (status int, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigIPWhiteIsOpen).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config.Name = ConfigIPWhiteIsOpen
		config.Value = "2"
		err = db.DB.Create(&config).Error
	}
	status = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetIPWhiteIsOpen(status int) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigIPWhiteIsOpen).Update("value", status).Error
	return
}

func (r *configRepo) GetCmsSiteInfo() (configs []model.Config, err error) {

	err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigCmsSiteName, ConfigCmsLoginIcon, ConfigCmsLoginBackend, ConfigCmsPageIcon, ConfigCmsMenuIcon}).Find(&configs).Error
	return
}

func (r *configRepo) SetCmsSiteName(value string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigCmsSiteName).Update("value", value).Error
	return
}

func (r *configRepo) SetCmsLoginIcon(value string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigCmsLoginIcon).Update("value", value).Error
	return
}

func (r *configRepo) SetCmsLoginBackend(value string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigCmsLoginBackend).Update("value", value).Error
	return
}

func (r *configRepo) SetCmsPageIcon(value string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigCmsPageIcon).Update("value", value).Error
	return
}

func (r *configRepo) SetCmsMenuIcon(value string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigCmsMenuIcon).Update("value", value).Error
	return
}

func (r *configRepo) GetJPushAuthInfo() (configs []model.Config, err error) {

	err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigJPushAppKey, ConfigJPushMasterSecret}).Find(&configs).Error
	return
}

func (r *configRepo) SetJPushAuthInfo(appKey, masterSecret string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigJPushAppKey).Update("value", appKey).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigJPushMasterSecret).Update("value", masterSecret).Error
	return
}

func (r *configRepo) GetFeihuAuthInfo() (configs []model.Config, err error) {

	var count int64
	err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigFeihuAppKey, ConfigFeihuAppSecret}).Find(&configs).Count(&count).Error
	if count != 2 {
		err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigFeihuAppKey, ConfigFeihuAppSecret}).Delete(&configs).Error
		err = db.DB.Model(model.Config{}).Create(&model.Config{Name: ConfigFeihuAppKey, Value: "d6KFWZ"}).Error
		err = db.DB.Model(model.Config{}).Create(&model.Config{Name: ConfigFeihuAppSecret, Value: "j0DCfAbUbLZw0q6E"}).Error
		return r.GetFeihuAuthInfo()
	}
	return
}

func (r *configRepo) SetFeihuAuthInfo(appKey, appSecret string) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigFeihuAppKey).Update("value", appKey).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigFeihuAppSecret).Update("value", appSecret).Error
	return
}

func (r *configRepo) GetDepositInfo() (configs []model.Config, err error) {

	var count int64
	err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigDepositHTML, ConfigDepositSWITCH, ConfigDepositURL}).Find(&configs).Count(&count).Error
	if count != 3 {
		err = db.DB.Model(model.Config{}).Create(&model.Config{Name: ConfigDepositHTML, Value: "<html><body><p>暂未开放，充值请联系接待！</p></body></html>"}).Error
		err = db.DB.Model(model.Config{}).Create(&model.Config{Name: ConfigDepositURL, Value: ""}).Error
		err = db.DB.Model(model.Config{}).Create(&model.Config{Name: ConfigDepositSWITCH, Value: "1"}).Error
		err = db.DB.Model(model.Config{}).Where("name IN ?", []string{ConfigDepositHTML, ConfigDepositSWITCH, ConfigDepositURL}).Find(&configs).Count(&count).Error
	}
	return
}

func (r *configRepo) SetDepositInfo(html, url string, open int) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDepositHTML).Update("value", html).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDepositURL).Update("value", url).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDepositSWITCH).Update("value", open).Error
	return
}

func (r *configRepo) GetRegisterConfig() (config model.SettingConfig, err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigRegisterConfig).First(&config).Error
	return
}

func (r *configRepo) GetSystemConfig() (p *model.ParameterConfigResp, err error) {
	defaultConfigByte, _ := util.MarshalJSONByDefault(&model.ParameterConfigResp{}, true)
	defaultConfig := string(defaultConfigByte.([]byte))
	m := model.ParameterConfigResp{}
	config := model.SettingConfig{}
	if err = db.DB.Where("config_type = ?", ConfigParameterConfig).First(&config).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if config.Content != "" {
		defaultConfig = config.Content
	}
	if err = json.Unmarshal([]byte(defaultConfig), &m); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
		return nil, err
	}
	cf, err := util.MarshalJSONByDefault(&m, false)
	if err != nil {
		return nil, err
	}
	m = *cf.(*model.ParameterConfigResp)

	return &m, nil
}

func (r *configRepo) UpdateParameterConfig(settingConfig *model.SettingConfig) (err error) {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"content"}),
	}).Create(&settingConfig).Error
}

func (r *configRepo) GetWithdrawConfig(lang string) (p *model.WithdrawConfigResp, err error) {
	Config_cn := `{
		"min": 100,
		"max": 50000,
		"columns": [
			{
				"name": "提款人:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 4
			},
			{
				"name": "银行:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 3
			},
			{
				"name": "卡号:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 2
			},
			{
				"name": "提现金额",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 1
			}
		]
	}`

	Config_en := `{
		"min": 100,
		"max": 50000,
		"columns": [
			{
				"name": "withdrawer:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 4
			},
			{
				"name": "bank:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 3
			},
			{
				"name": "card number:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 2
			},
			{
				"name": "withdrawal Amount",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 1
			}
		]
	}`

	Config_ja := `{
		"min": 100,
		"max": 50000,
		"columns": [
			{
				"name": "引き出し者:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 4
			},
			{
				"name": "銀行:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 3
			},
			{
				"name": "カード番号:",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 2
			},
			{
				"name": "出金額",
				"default_content": "",
				"required": 1,
				"default_content_modify": 1,
				"sort": 1
			}
		]
	}`
	var default_Config string
	switch lang {
	case "en_US":
		default_Config = Config_en
	case "zh_CN":
		default_Config = Config_cn
	case "ja":
		default_Config = Config_ja
	default:
		default_Config = Config_cn
	}

	m := model.WithdrawConfigResp{}
	m_c := model.WithdrawConfigResp{}
	if err = json.Unmarshal([]byte(default_Config), &m_c); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
		return nil, err
	}
	config := model.SettingConfig{}
	if err = db.DB.Where("config_type = ?", ConfigWithdrawConfig).First(&config).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	if err = json.Unmarshal([]byte(config.Content), &m); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
		return nil, err
	}
	m_c.Max = m.Max
	m_c.Min = m.Min
	sort.Sort(m_c.Columns)
	return &m_c, nil
}

func (r *configRepo) UpdateWithdrawConfig(settingConfig *model.SettingConfig) (err error) {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"content"}),
	}).Create(&settingConfig).Error
}

func (r *configRepo) GetLoginConfig() (config model.SettingConfig, err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigLoginConfig).First(&config).Error
	return
}

func (r *configRepo) UpdateRegisterConfig(config string) (err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigRegisterConfig).Update("content", config).Error
	return
}

func (r *configRepo) UpdateLoginConfig(config string) (err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigLoginConfig).Update("content", config).Error
	return
}

func (r *configRepo) GetAnnouncementConfig() (config model.SettingConfig, err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigAnnouncementConfig).First(&config).Error
	return
}

func (r *configRepo) UpdateAnnouncementConfig(config string) (err error) {

	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigAnnouncementConfig).Update("content", config).Error
	return
}

func (r *configRepo) UpdateSignConfig(settingConfig *model.SettingConfig) error {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"content"}),
	}).Create(&settingConfig).Error
}

func (r *configRepo) GetSignConfig() (config model.SettingConfig, err error) {
	err = db.DB.Table("setting_configs").Where("config_type = ?", ConfigSignConfig).First(&config).Error
	return
}

func (r *configRepo) GetUserAgreement() (content string, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigUserAgreement).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		content = DefaultUserAgreement
		config.Name = ConfigUserAgreement
		config.Value = content
		err = db.DB.Create(&config).Error
		return
	}
	content = config.Value
	return
}

func (r *configRepo) SetUserAgreement(content string) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigUserAgreement).Update("value", content).Error
	return
}

func (r *configRepo) GetPrivacyPolicy() (content string, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigPrivacyPolicy).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		content = DefaultPrivacyPolicy
		config.Name = ConfigPrivacyPolicy
		config.Value = content
		err = db.DB.Create(&config).Error
		return
	}
	content = config.Value
	return
}

func (r *configRepo) SetPrivacyPolicy(content string) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigPrivacyPolicy).Update("value", content).Error
	return
}

func (r *configRepo) GetAboutUs() (content string, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigAboutUs).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		content = DefaultAboutUs
		config.Name = ConfigAboutUs
		config.Value = content
		err = db.DB.Create(&config).Error
		return
	}
	content = config.Value
	return
}

func (r *configRepo) SetAboutUs(content string) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigAboutUs).Update("value", content).Error
	return
}

func (r *configRepo) GetPrivilegeUserFreezeIsOpen() (status int, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigPrivilegeUserFreeze).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config.Name = ConfigPrivilegeUserFreeze
		config.Value = "2"
		err = db.DB.Create(&config).Error
	}
	if err == nil {
		status = util.String2Int(config.Value)
	}
	return
}

func (r *configRepo) SetPrivilegeUserFreezeIsOpen(status int) (err error) {

	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigPrivilegeUserFreeze).Update("value", status).Error
	return
}

func (r *configRepo) GetDefaultGroupIndex() (index int, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultGroupIndex).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config.Name = ConfigDefaultGroupIndex
		config.Value = "0"
		err = db.DB.Create(&config).Error
	}
	index = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetDefaultGroupIndex(index int) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultGroupIndex).Update("value", index).Error
	return
}

func (r *configRepo) GetDefaultFriendIndex() (index int, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultFriendIndex).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config.Name = ConfigDefaultFriendIndex
		config.Value = "0"
		err = db.DB.Create(&config).Error
	}
	index = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetDefaultFriendIndex(index int) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultFriendIndex).Update("value", index).Error
	return
}

func (r *configRepo) GetDefaultIsOpen() (index int, err error) {
	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultIsOpen).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config.Name = ConfigDefaultIsOpen
		config.Value = "2"
		err = db.DB.Create(&config).Error
	}
	index = util.String2Int(config.Value)
	return
}

func (r *configRepo) SetDefaultIsOpen(index int) (err error) {
	err = db.DB.Model(model.Config{}).Where("name = ?", ConfigDefaultIsOpen).Update("value", index).Error
	return
}
