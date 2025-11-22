	create table EA_Tags
			(
			  id int identity(1,1) primary key,
			  name_ar nvarchar(250),
			  name_en nvarchar(250),
			  description nvarchar(350)
			)
			
	create table EA_Tags_Dimentions
			(
			  id int identity(1,1),
			  ea_tag_id int not null references EA_tags(id),
			  object_type_id int not null  references objecttype (objecttypeid)
			)

			iNSERT INTO EA_Tags (name_en, name_ar, description)
VALUES 
('Business Dimension', N'البعد التجاري', N'يتعلق بالعمليات والأهداف التجارية للمؤسسة'),
('Technology Dimension', N'البعد التقني', N'يركز على البنية التحتية التقنية والتطبيقات والأنظمة'),
('Data Dimension', N'بعد البيانات', N'يعالج إدارة البيانات وجودتها وحوكمتها'),
('Organization Dimension', N'البعد التنظيمي', N'يصف هيكل المؤسسة والأدوار والمسؤوليات'),
('Strategy Dimension', N'البعد الاستراتيجي', N'يتناول الرؤية والأهداف والخطط طويلة المدى');