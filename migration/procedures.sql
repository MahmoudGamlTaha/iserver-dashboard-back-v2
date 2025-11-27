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


