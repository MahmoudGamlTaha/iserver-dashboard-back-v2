USE [iserver-light]
GO

/****** Object:  StoredProcedure [import].[InsertOrUpdateImportedObject]    Script Date: 11/24/2025 11:03:23 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE OR ALTER PROCEDURE [import].[InsertOrUpdateImportedObject]
    @ObjectID               UNIQUEIDENTIFIER,
    @VersionID              UNIQUEIDENTIFIER,
    @GeneralTypeID          INT,
    @ObjectGroupID          UNIQUEIDENTIFIER,
    @ObjectName             NVARCHAR(1024),
    @ObjectDescription      NVARCHAR(MAX),
    @RichTextDescription    NVARCHAR(MAX),
    @UserID                 INT,
    @ProfileID              INT,
    @VersionName            NVARCHAR(128),
    @LibraryId              UNIQUEIDENTIFIER,
    @ParentObjectID         UNIQUEIDENTIFIER = NULL,
    @ObjectTypeID           INT,
    @Lock                   BIT,
    @IgnoreCheckout         BIT,
    @RestoreToImportFolder  BIT,
    @HasVisioAlias          BIT,
    @VisioAlias             NVARCHAR(MAX),
    @Result                 BIT = 0 OUTPUT,
    @Error                  NVARCHAR(1024) = '' OUTPUT
AS

BEGIN

    SET NOCOUNT ON       
    
    SET @Result = 0
    SET @Error = ''

    BEGIN TRAN

        IF NOT EXISTS (SELECT 1 FROM Object WHERE ObjectId = @ObjectId)
        BEGIN
            DECLARE @CreateSpec [object].udttCreateComponentSpec
            INSERT INTO @CreateSpec(ObjectId, VersionId, ObjectTypeID, ObjectName, ObjectDescription, ParentObjectId, ObjectGroupId, LibraryId, CreateContentsRecord, ExactObjectTypeID, VisioAlias, HasVisioAlias, UserVersionNo, RichTextDescription, ApprovalStatus)
            VALUES(@ObjectID, @VersionId, @GeneralTypeID, @ObjectName, @ObjectDescription, @ParentObjectID, @ObjectGroupID, @LibraryId, 1, @ObjectTypeID, @VisioAlias, @HasVisioAlias, @VersionName, @RichTextDescription, dbo.const_ApprovalStatus_Approved())

			DECLARE @return_status INT
			
            EXEC [object].uspComponentInsert @CreateSpec, @UserId, @ProfileId, @return_status

			IF @return_status <> 0 RETURN @return_status
			
            -- Set it's imported flag to true:
            UPDATE	dbo.[Object]
            SET		IsImported = 1
            WHERE	ObjectID = @ObjectID
                
            -- Insert Object Automatically checks out - we want object checked in
            UPDATE  [Version]
            SET     IsCheckedOut=0, 
                    CheckedOutBy=NULL,
                    DateModified = GETDATE(),
                    ModifiedBy = @UserID
            WHERE   ID = @VersionID

            UPDATE  [Object]
            SET     IsCheckedOut = 0,
                    CheckedOutUserId = NULL,
                    CheckedInVersionId = CurrentVersionId
            WHERE   [Object].ObjectID = @ObjectId

            IF @Lock = 1 
            BEGIN           
                UPDATE dbo.[Object]
                SET    Locked = 1
                WHERE  ObjectID = @ObjectID
            END

            SET @Result = 1
        END
        ELSE
        BEGIN
            -- This should be handled in the UI, but last resort check here
            IF EXISTS (SELECT 1 FROM ObjectPermissions AS op WHERE op.ObjectID = @ObjectID AND op.ProfileID = @ProfileID AND op.HasRead = CAST(1 AS BIT))
            BEGIN
                IF @IgnoreCheckout = 1 OR EXISTS (SELECT 1 FROM Version WHERE ID = @VersionID AND (IsCheckedOut != CAST(1 AS BIT) OR CheckedOutBy = @UserId))
                BEGIN
                    --Update object with imported details
                    UPDATE	o
                    SET     o.ObjectName = @ObjectName,
                            o.ObjectDescription = @ObjectDescription,
                            o.RichTextDescription = @RichTextDescription,
                            o.DateModified = GETDATE(),
                            o.ModifiedBy = @UserID,
                            o.VisioAlias = @VisioAlias, 
                            o.HasVisioAlias = @HasVisioAlias
                    FROM	Object AS o
                    WHERE	o.ObjectID = @ObjectID

                    UPDATE  [Version]
                    SET     ObjectName = @ObjectName,
                            ObjectDescription = @ObjectDescription,
                            RichTextDescription = @RichTextDescription,
                            UserVersionNo = @VersionName,
                            HasVisioAlias = @HasVisioAlias,
                            VisioAlias = @VisioAlias
                    WHERE   ID = @VersionID

                    -- Assign to group if needed
                    IF @ObjectGroupID IS NOT NULL AND NOT EXISTS (SELECT 1 FROM ObjectGroupContents WHERE ObjectGroupId = @ObjectGroupId AND ObjectId = @ObjectId)
                    BEGIN
                        INSERT INTO ObjectGroupContents (ObjectGroupId, ObjectId)
                        VALUES (@ObjectGroupId, @ObjectId)						
                    END

                    --If this option is on, put the shape back
                    --into the folder it was imported into
                    IF @RestoreToImportFolder = 1 AND @ParentObjectID IS NOT NULL
                    BEGIN
                        DECLARE @DocumentObjectsId INT
                        
                        SELECT TOP 1 @DocumentObjectsId = ID 
                        FROM	vwCurrentObjectContents (NOEXPAND)
                        WHERE	DocumentObjectID= @ParentObjectID
                        AND		ObjectID = @ObjectID 
                        
                        IF @DocumentObjectsId IS NULL
                        BEGIN
                            INSERT INTO ObjectContents
                            (
                                    [DocumentObjectID],
                                    [ContainerVersionID],
                                    [ObjectID],
                                    [ContainmentType],
                                    [Instances],
                                    [IsShortCut],
                                    [DateCreated],
                                    [CreatedBy],
                                    [DateModified],
                                    [ModifiedBy]
                            )
                            SELECT	@ParentObjectID,
                                    container.CurrentVersionId,
                                    @ObjectID,
                                    dbo.const_ContainmentType_ShapeInFolder(),
                                    1,
                                    1,
                                    GETDATE(),
                                    @UserID,
                                    GETDATE(),
                                    @UserID
                            FROM [Object] AS container
                            WHERE container.ObjectID = @ParentObjectID
                        END
                    END

                    -- Explicitly lock if not checked out
                    IF (@Lock = 1 AND EXISTS (SELECT 1 FROM [Version] WHERE ID = @VersionID AND IsCheckedOut != CAST(1 AS BIT)))  
                    BEGIN           
                        UPDATE dbo.[Object]
                        SET    Locked = 1
                        WHERE  ObjectID = @ObjectID
                    END

                    SET @Result = 1
                END
                ELSE
                BEGIN
                    SET @Error = 'Unable to update object because it is checked-out'
                END
            END
            ELSE
            BEGIN
                SET @Error = 'Unable to update object because you do not have permission'
            END

        END

    COMMIT

END
GO
------ 

USE [iserver-light]
GO

/****** Object:  StoredProcedure [dbo].[usp_InsertNewVersionForExistingObject]    Script Date: 11/30/2025 10:36:16 AM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[usp_InsertNewVersionForExistingObject]         
    @ObjectId UNIQUEIDENTIFIER,
    @NewVersionId UNIQUEIDENTIFIER,
    @UserVersionNo NVARCHAR(128),
    @UserId INT,
    @NewVersionIsCheckedOut BIT = 1
AS

DECLARE @CurrentVersionId UNIQUEIDENTIFIER
DECLARE @SystemVersionNo INT
DECLARE @ObjectTypeId INT

SET NOCOUNT ON

SELECT	@CurrentVersionId = CurrentVersionId,
        @ObjectTypeId = ObjectTypeId
FROM	Object 
WHERE	ObjectId = @ObjectId

SELECT	@SystemVersionNo = SystemVersionNo + 1 
FROM	Version 
WHERE	ObjectId = @ObjectId 
AND		[Id] = @CurrentVersionId

BEGIN TRY  
BEGIN TRAN          
            

    INSERT INTO	Version 
    (
                    [Id], 
                    ObjectId, 
                    ObjectName,
                    ObjectDescription,
                    RichTextDescription,
                    UserVersionNo,
                    SystemVersionNo, 
                    HasVisioAlias,
                    VisioAlias,
                    VersionImage,
                    IsCheckedOut, 
                    CheckedOutBy,
                    CheckInReason,
                    FileExtension,
                    PrefixForRollback,
                    SuffixForRollback,
                    ProvenanceIdForRollback,
                    ProvenanceVersionIdForRollback,
                    RequiresShapeSheetUpdateForRollback,
                    HasValidContentsHistory,
                    HasValidVisioPageInstances,
                    ApprovalStatus,
                    DateCreated,
                    CreatedBy,
                    DateModified,
                    ModifiedBy 
    )          
    SELECT			@NewVersionId, 
                    ObjectId, 
                    ObjectName,
                    ObjectDescription,
                    RichTextDescription,
                    UserVersionNo = @UserVersionNo, 
                    SystemVersionNo = @SystemVersionNo, 
                    HasVisioAlias,
                    VisioAlias,
                    VersionImage, 
                    IsCheckedOut = @NewVersionIsCheckedOut,
                    CheckedOutBy = (CASE WHEN @NewVersionIsCheckedOut = 1 THEN @UserId ELSE NULL END), 
                    CheckInReason = '',
                    FileExtension,
                    PrefixForRollback,
                    SuffixForRollback,
                    ProvenanceIdForRollback,
                    ProvenanceVersionIdForRollback,
                    RequiresShapeSheetUpdateForRollback,
                    1, -- New versions always have valid history
                    HasValidVisioPageInstances,
                    dbo.const_ApprovalStatus_Approved(),
                    GETDATE(),
                    @UserId,
                    GETDATE(),
                    @UserId 
    FROM			Version
    WHERE			Version.[Id] = @CurrentVersionId 
    AND				Version.ObjectId = @ObjectId 


    UPDATE Object
    SET IsCheckedOut = @NewVersionIsCheckedOut,
    CheckedOutUserId = (CASE WHEN @NewVersionIsCheckedOut = 1 THEN @UserId ELSE NULL END)
    WHERE ObjectId = @ObjectId



    IF @ObjectTypeId = dbo.const_GeneralType_VisioShape()
        UPDATE	Object
        SET		CurrentVersionId = @NewVersionId, 
                CheckedInVersionId = CASE WHEN @NewVersionIsCheckedOut = 0 THEN @NewVersionId ELSE CheckedInVersionId END
        WHERE	ObjectId = @ObjectId
    ELSE
        UPDATE	Object
        SET		CurrentVersionId = @NewVersionId, 
                CheckedInVersionId = CASE WHEN @NewVersionIsCheckedOut = 0 THEN @NewVersionId ELSE CheckedInVersionId END
        WHERE	ObjectId = @ObjectId

    -- Insert new ObjectContents Versions
    INSERT INTO ObjectContents(DocumentObjectID, ContainerVersionID, ObjectID, ContainmentType, Instances, IsShortCut, DateCreated, CreatedBy, DateModified, ModifiedBy)
    SELECT DocumentObjectID, @NewVersionId, ObjectID, ContainmentType, Instances, IsShortcut, GETDATE(), @UserId, GETDATE(), @UserId
    FROM ObjectContents AS oc
    WHERE DocumentObjectID = @ObjectId AND oc.ContainerVersionID = @CurrentVersionId

    -- Clone Any RelationDocument entries
    INSERT INTO RelationDocument(RelationshipId, DocumentId, DocumentVersionId, Instances, ShapeSheetKeysRequiringUpdateId)
    SELECT rd.RelationshipId, rd.DocumentId, @NewVersionId, Instances, rd.ShapeSheetKeysRequiringUpdateId
    FROM RelationDocument AS rd
    WHERE rd.DocumentId = @ObjectId AND rd.DocumentVersionId = @CurrentVersionId
    
    -- Clone any VisioPageShapes entries
    INSERT INTO VisioPageShapes(DocumentObjectId, DocumentVersionId, ShapeObjectId, VisioPageId, VisioShapeId, VisioShapeHeight, VisioShapeWidth, VisioShapeLeft, VisioShapeTop, VisioShapeRight, VisioShapeBottom, SequenceFlow)
    SELECT vps.DocumentObjectId, @NewVersionId, ShapeObjectId, VisioPageId, VisioShapeId, VisioShapeHeight, VisioShapeWidth, VisioShapeLeft, VisioShapeTop, VisioShapeRight, VisioShapeBottom, SequenceFlow
    FROM VisioPageShapes AS vps
    WHERE DocumentObjectId = @ObjectId AND DocumentVersionId = @CurrentVersionId

    -- Clone any VisioPageRelationships entries
    INSERT INTO VisioPageRelationships(DocumentObjectId, DocumentVersionId, RelationshipId, VisioPageId, BeginVisioId, EndVisioId, ConnectorVisioId, ConnectorMasterBaseId, ConnectorBeginX, ConnectorBeginY, ConnectorEndX, ConnectorEndY)
    SELECT vpr.DocumentObjectId, @NewVersionId, RelationshipId, VisioPageId, BeginVisioId, EndVisioId, ConnectorVisioId, ConnectorMasterBaseId, ConnectorBeginX, ConnectorBeginY, ConnectorEndX, ConnectorEndY
    FROM VisioPageRelationships AS vpr
    WHERE DocumentObjectId = @ObjectId AND DocumentVersionId = @CurrentVersionId

    -- Insert new attribute value versions
    INSERT INTO AttributeValue (ObjectId, VersionId, AttributeId, DataType, ValueBigInt, ValueDate, ValueFloat, ValueText, ValueRichText, ValueHTML, DateCreated, CreatedBy, DateModified, ModifiedBy)
    SELECT ObjectId, 
           @NewVersionId, 
           AttributeId, 
           DataType,
           ValueBigInt,
           ValueDate,
           ValueFloat,
           ValueText, 
           ValueRichText,
		   ValueHTML,
           DateCreated, 
           CreatedBy, 
           DateModified, 
           ModifiedBy
    FROM AttributeValue 
    WHERE ObjectId = @ObjectId 
        AND VersionId = @CurrentVersionId

COMMIT TRAN	

END TRY	
    
BEGIN CATCH
    ROLLBACK TRAN
END CATCH

SELECT @NewVersionId AS VersionID,  @ObjectTypeId AS ObjectTypeID

RETURN 1
GO




